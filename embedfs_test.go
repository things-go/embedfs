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

func performGetRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
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
	w = performGetRequest(r, "/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/index")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/hello")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/ioo")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/ixx")
	require.Equal(t, http.StatusOK, w.Code)

	// 系统文件
	w = performGetRequest(r, "/xss/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/xmg/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/xss")
	require.Equal(t, http.StatusMovedPermanently, w.Code)
	w = performGetRequest(r, "/xmg")
	require.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performGetRequest(r, "/5.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/6.jpg")
	require.Equal(t, http.StatusOK, w.Code)

	// 嵌入绑定的文件
	w = performGetRequest(r, "/css/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/img/")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/css")
	require.Equal(t, http.StatusMovedPermanently, w.Code)
	w = performGetRequest(r, "/img")
	require.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performGetRequest(r, "/1.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/4.png")
	require.Equal(t, http.StatusOK, w.Code)
	w = performGetRequest(r, "/xxxx.png")
	require.Equal(t, http.StatusOK, w.Code)
}
