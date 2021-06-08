package testdata

import "embed"

//go:embed static/css static/img static/views static/1.png static/index.html
var Staticfs embed.FS
