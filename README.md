# embedfs

[![GoDoc](https://godoc.org/github.com/things-go/embedfs?status.svg)](https://godoc.org/github.com/things-go/embedfs)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/embedfs?tab=doc)
[![Build Status](https://www.travis-ci.com/things-go/embedfs.svg?branch=master)](https://www.travis-ci.com/things-go/embedfs)
[![codecov](https://codecov.io/gh/things-go/embedfs/branch/master/graph/badge.svg)](https://codecov.io/gh/things-go/embedfs)
![Action Status](https://github.com/things-go/embedfs/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/embedfs)](https://goreportcard.com/report/github.com/things-go/embedfs)
[![License](https://img.shields.io/github/license/things-go/embedfs)](https://github.com/things-go/embedfs/raw/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/embedfs)](https://github.com/things-go/embedfs/tags)


## Installation

```bash
    go get github.com/things-go/embedfs
```

## Example

[embedmd]:# (_example/main.go go)
```go
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
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
