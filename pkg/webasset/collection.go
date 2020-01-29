package webasset

import (
	"errors"
	"html/template"
	"io"
	"strings"
)

var (
	ErrTemplateNotFound = errors.New("template not found")

	htmlTemplatesSuffix = ".html"
	stylesheetsSuffix   = ".css"
	scriptsSuffix       = ".js"
)

type MustAssetFunc func(string) []byte

type AssetCollection struct {
	templates map[string]*template.Template

	stylesheet map[string][]byte
	dependency map[string][]byte
	script     map[string][]byte

	renderable struct {
		Stylesheets  []string
		Dependencies []string
		Scripts      []string
	}
}

type renderable struct {
	Stylesheets  []string
	Dependencies []string
	Scripts      []string
}

func (a *AssetCollection) mustParseTemplate(name string, mustAsset MustAssetFunc) {
	t, err := template.New(name).Parse(string(mustAsset(name)))
	if err != nil {
		panic(err)
	}
	a.templates[name] = t
}

func (a *AssetCollection) AddStylesheet(name string, data []byte) {
	_, ok := a.stylesheet[name]
	a.stylesheet[name] = data

	if !ok {
		a.renderable.Stylesheets = append(a.renderable.Stylesheets, name)
	}
}

func (a *AssetCollection) AddDependency(name string, data []byte) {
	_, ok := a.dependency[name]
	a.dependency[name] = data

	if !ok {
		a.renderable.Dependencies = append(a.renderable.Dependencies, name)
	}
}

func (a *AssetCollection) AddScript(name string, data []byte) {
	_, ok := a.script[name]
	a.script[name] = data

	if !ok {
		a.renderable.Scripts = append(a.renderable.Scripts, name)
	}
}

func (a *AssetCollection) ExecuteTemplate(wr io.Writer, name string) error {
	t, ok := a.templates[name]
	if !ok {
		return ErrTemplateNotFound
	}
	return t.ExecuteTemplate(wr, name, a.renderable)
}

func MustParseCollection(assetNames []string, mustAsset MustAssetFunc) *AssetCollection {
	collection := &AssetCollection{
		templates:  make(map[string]*template.Template),
		dependency: make(map[string][]byte),
		stylesheet: make(map[string][]byte),
		script:     make(map[string][]byte),
		renderable: renderable{
			Stylesheets:  []string{},
			Dependencies: []string{},
			Scripts:      []string{},
		},
	}
	for _, name := range assetNames {
		if strings.HasSuffix(name, htmlTemplatesSuffix) {
			collection.mustParseTemplate(name, mustAsset)
		} else if strings.HasSuffix(name, stylesheetsSuffix) {
			collection.AddStylesheet(name, mustAsset(name))
		} else if strings.HasSuffix(name, scriptsSuffix) {
			if strings.Contains(name, "/") {
				collection.AddDependency(name, mustAsset(name))
			} else {
				collection.AddScript(name, mustAsset(name))
			}
		}
	}
	return collection
}
