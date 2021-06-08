package embedfs

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/things-go/embedfs/testdata"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStaticFileFs(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 静态html文件
	HTML(r, WWW{
		EmbedFs:          testdata.Staticfs,
		EmbedTplPatterns: []string{"static/index.html", "static/views/*"},
		TplPatterns:      []string{"testdata/static/ixx.html"},
		RelativePathToTpl: map[string]string{
			"/":      "index.html",
			"/index": "index.html",
			"/hello": "hello.html",
			"/ioo":   "ioo.html",
			"/ixx":   "ixx.html",
		},
	})

	// 系统文件
	r.StaticFS("/xss", http.Dir("testdata/static/css"))
	r.StaticFS("/xmg", http.Dir("testdata/static/img"))
	r.StaticFile("/5.png", "testdata/static/5.png")
	r.StaticFile("/6.jpg", "testdata/static/views/6.jpg")
	// 嵌入绑定的文件
	cssFs, _ := fs.Sub(testdata.Staticfs, "static/css")
	r.StaticFS("/css", http.FS(cssFs))
	imgFs, _ := fs.Sub(testdata.Staticfs, "static/css")
	r.StaticFS("/img", http.FS(imgFs))
	StaticFileFS(r, "/1.png", "static/1.png", http.FS(testdata.Staticfs))
	StaticFileFS(r, "/4.png", "static/views/4.png", http.FS(testdata.Staticfs))
	StaticFileFS(r, "/xxxx.png", "x.png", http.FS(testdata.Staticfs))

	var w *httptest.ResponseRecorder

	require.Panics(t, func() {
		StaticFileFS(r, "/*.png", "static/1.png", http.FS(testdata.Staticfs))
	})

	// html
	w = performRequest(r, http.MethodGet, "/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/index")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/hello")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/ioo")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/ixx")
	require.Equal(t, http.StatusOK, w.Code)

	// 系统文件
	w = performRequest(r, http.MethodGet, "/xss/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/xmg/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/xss")
	require.Equal(t, http.StatusMovedPermanently, w.Code)
	w = performRequest(r, http.MethodGet, "/xmg")
	require.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(r, http.MethodGet, "/5.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/6.jpg")
	require.Equal(t, http.StatusOK, w.Code)

	// 嵌入绑定的文件
	w = performRequest(r, http.MethodGet, "/css/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/img/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/css")
	require.Equal(t, http.StatusMovedPermanently, w.Code)
	w = performRequest(r, http.MethodGet, "/img")
	require.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(r, http.MethodGet, "/1.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/4.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performRequest(r, http.MethodGet, "/xxxx.png")
	require.Equal(t, http.StatusOK, w.Code)
}
