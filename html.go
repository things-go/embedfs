package embedfs

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WWW HTML config
type WWW struct {
	EmbedFs           fs.FS             // 如果不为nil使用embed fs
	EmbedTplPatterns  []string          // embed fs模板patterns, see: template的ParseFS和ParseGlob
	TplPatterns       []string          // 文件模板patterns, see: template的ParseFS和ParseGlob
	RelativePathToTpl map[string]string // url相对路径对模板名映射
}

// HTML 静态页面处理
// 只可以设一次静态文件
func HTML(engine *gin.Engine, c WWW) *gin.Engine {
	tpl := template.New("www")
	if c.EmbedFs != nil {
		tpl = template.Must(tpl.ParseFS(c.EmbedFs, c.EmbedTplPatterns...))
	}
	for _, pattern := range c.TplPatterns {
		tpl = template.Must(tpl.ParseGlob(pattern))
	}
	engine.SetHTMLTemplate(tpl)
	// 静态文件URL
	for _path, name := range c.RelativePathToTpl {
		tmpPath, tmpName := _path, name
		engine.GET(tmpPath, func(c *gin.Context) {
			c.HTML(http.StatusOK, tmpName, nil)
		})
	}
	return engine
}
