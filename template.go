package gow

import (
	"html/template"
	"path/filepath"
	"github.com/labstack/gommon/log"
)

type TemplateEngine struct {
	templatesDir    string
	templatesSuffix string
	Cached          bool
	IsInit          bool
	templates       map[string]*template.Template
}

func (tple *TemplateEngine) CreateTemplate(templateName []string, data map[string]interface{}) {
	tpl := template.New(templateName[0]).Funcs(data)
	if data != nil {
		tpl.Funcs(data)
	}
	tpl.ParseFiles()
}

func (tple *TemplateEngine) Render(tpl string, data map[string]interface{}) (b []byte, err error) {
	return nil, nil
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
		if tple.templates == nil {
			tple.templates = make(map[string]*template.Template)
		}
		tple.templatesDir = "templates/"
		tple.templatesSuffix = ".html"
		
		layouts, err := filepath.Glob(tple.templatesDir + "layouts/*" + tple.templatesSuffix)
		if err != nil {
			log.Error(err)
		}
		
		includes, err := filepath.Glob(tple.templatesDir + "includes/*" + tple.templatesSuffix)
		if err != nil {
			log.Fatal(err)
		}
		
		for _, layout := range layouts {
			files := append(includes, layout)
			tple.templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
		}
		tple.IsInit = true
	}
}
