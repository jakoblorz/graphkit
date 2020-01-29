package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/jakoblorz/graphkit/pkg/mime"
	"github.com/jakoblorz/graphkit/pkg/webasset"
	"golang.org/x/net/websocket"
)

//go:generate go-bindata -debug -prefix "../../assets" -o assets.go ../../assets/...

var (
	wellKnownFileEndings = map[string]string{
		".css": "text/css",
		".js":  "application/javascript",
	}
	webBasedRenderingConfiguration = []string{
		"index.html",
		"index.css",

		"lib/babel.min.js",
		"lib/viz.js",

		"receive.js",
		"render.js",
	}
	webVisualInit                  = `digraph {}`
	dotBasedRenderingConfiguration = []string{
		"index.html",
		"index.css",

		"lib/babel.min.js",
		"lib/react.development.js",
		"lib/react-dom.development.js",

		"receive.js",
		"react-app.js",
	}
	dotVisualInit = `<svg
		xmlns="http://www.w3.org/2000/svg" 
		xmlns:xlink="http://www.w3.org/1999/xlink"
	></svg>`
)

func serveTemplate(name string, loader *webasset.AssetCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		err := loader.ExecuteTemplate(w, name)
		if err != nil {
			panic(err)
		}
	}
}

func serveAsset(name string) http.HandlerFunc {
	asset := MustAsset(name)
	contentType := mime.DetectContentType(name, asset[:512], wellKnownFileEndings)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", contentType)
		io.Copy(w, bytes.NewReader(asset))
	}
}

func serveWebsocket(updateCh chan string, init string) websocket.Handler {

	sockets := []chan string{}
	mx := &sync.Mutex{}

	go func() {
		for {
			update, ok := <-updateCh
			if !ok {
				mx.Lock()
				for _, s := range sockets {
					close(s)
				}
				mx.Unlock()
				return
			}

			mx.Lock()
			for _, s := range sockets {
				s <- update
			}
			init = update
			mx.Unlock()
		}
	}()

	return func(ws *websocket.Conn) {
		ch := make(chan string, 1)

		mx.Lock()
		ch <- init
		sockets = append(sockets, ch)
		mx.Unlock()

		for {
			update, ok := <-ch
			if !ok {
				break
			}

			_, err := ws.Write([]byte(update))
			if err != nil {
				log.Printf("error on websocket write, closing: %+v\n", err)
				break
			}
		}

		mx.Lock()
		for i, cmp := range sockets {
			if cmp != ch {
				continue
			}

			sockets[i], sockets[len(sockets)-1] = sockets[len(sockets)-1], sockets[i]
			sockets = sockets[:len(sockets)-1]
			break
		}
		mx.Unlock()
	}
}

func main() {
	var (
		url = flag.String("url", ":8080", "specify the host:port combination")
		web = flag.Bool("web", false, "enable in-browser rendering of .dot file. Enable if graphviz is not installed")
	)
	flag.Parse()

	visualInit := dotVisualInit
	assetNamesToLoad := dotBasedRenderingConfiguration
	if *web {
		visualInit = webVisualInit
		assetNamesToLoad = webBasedRenderingConfiguration
	}

	loader := webasset.MustParseCollection(assetNamesToLoad, MustAsset)
	http.HandleFunc("/", serveTemplate("index.html", loader))
	for _, name := range assetNamesToLoad {
		path := fmt.Sprintf("/%s", name)
		if strings.HasSuffix(name, ".html") {
			http.HandleFunc(path, serveTemplate(name, loader))
		} else {
			http.HandleFunc(path, serveAsset(name))
		}
	}

	updateCh := make(chan string, 1)
	http.Handle("/ws", serveWebsocket(updateCh, visualInit))

	err := http.ListenAndServe(*url, nil)
	if err != nil {
		log.Printf("failed to serve application on %s: %+v\n", *url, err)
	}
}
