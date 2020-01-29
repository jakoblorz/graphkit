package mime

import (
	"net/http"
	"strings"
)

func DetectContentType(name string, buf []byte, wellKnownEndings map[string]string) string {
	if wellKnownEndings != nil {
		for suffix, contentType := range wellKnownEndings {
			if strings.HasSuffix(name, suffix) {
				return contentType
			}
		}
	}
	return http.DetectContentType(buf)
}
