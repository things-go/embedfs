package embedfs

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticFileFS(r gin.IRouter, relativePath, filepath string, fs http.FileSystem) gin.IRouter {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *gin.Context) {
		c.FileFromFS(filepath, fs)
	}
	r.GET(relativePath, handler)
	r.OPTIONS(relativePath, handler)
	return r
}
