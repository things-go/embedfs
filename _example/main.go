package main

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/things-go/embedfs"
	"github.com/things-go/embedfs/testdata"
)

func main() {
	r := gin.Default()

	// 静态html文件
	embedfs.HTML(r, embedfs.WWW{
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
	embedfs.StaticFileFS(r, "/1.png", "static/1.png", http.FS(testdata.Staticfs))
	embedfs.StaticFileFS(r, "/4.png", "static/views/4.png", http.FS(testdata.Staticfs))

	err := r.Run(":9000")
	if err != nil {
		panic(err)
	}
}
