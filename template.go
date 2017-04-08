package gow

import (
	"html/template"
	"path/filepath"
	"github.com/biezhi/gow/bpool"
	"net/http"
	"path"
	"bytes"
)

type TemplateEngine struct {
	templatesDir    string
	templatesSuffix string
	Cached          bool
	IsInit          bool
	templates       map[string]*template.Template
}

var bufpool *bpool.BufferPool

//func (tple *TemplateEngine) CreateTemplate(templateName string, data map[string]interface{}) template.Template {
//	tpl := template.New(templateName).Funcs(data)
//	if data != nil {
//		tpl.Funcs(data)
//	}
//
//	tpl.ParseFiles()
//}

func (tple *TemplateEngine) Render(w http.ResponseWriter, tplName string, data map[string]interface{}) (b []byte, err error) {
	
	tplPath := path.Join(tple.templatesDir, tplName + tple.templatesSuffix)
	tpl := template.New(tplName +  tple.templatesSuffix)
	template.Must(tpl.ParseFiles(tplPath))
	
	var buf bytes.Buffer
	e := tpl.Execute(&buf, data)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}

// close template cache
func (tple *TemplateEngine) CloseCache() {
	tple.Cached = false
	tple.templates = make(map[string]*template.Template)
}

func NewTemplateEngine() *TemplateEngine {
	tple := TemplateEngine{}
	tple.Init()
	return &tple
}

func (tple *TemplateEngine) Init() {
	if !tple.IsInit {
		logger.Info("Init Template Engine")
		if tple.templates == nil {
			tple.templates = make(map[string]*template.Template)
		}
		tple.templatesDir = "templates/"
		tple.templatesSuffix = ".html"
		
		layouts, err := filepath.Glob(tple.templatesDir + "layouts/*" + tple.templatesSuffix)
		if err != nil {
			logger.Error(err.Error())
		}
		
		includes, err := filepath.Glob(tple.templatesDir + "includes/*" + tple.templatesSuffix)
		if err != nil {
			logger.Error(err.Error())
		}
		
		for _, layout := range layouts {
			files := append(includes, layout)
			tple.templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
		}
		
		tple.IsInit = true
	}
}
