package gow

import (
	"html/template"
	"path/filepath"
	"path"
	"bytes"
	"github.com/biezhi/agon/utils"
)

type TemplateEngine struct {
	templatesDir    string
	templatesSuffix string
	Cached          bool
	IsInit          bool
	templates       map[string]*template.Template
}

const (
	COMM_TPL_COMMONS  = "commons"
	COMM_TPL_INCLUDES = "includes"
	COMM_TPL_LAYOUTS  = "layouts"
)

func (tple *TemplateEngine) CreateTemplate(tplName string) *template.Template {
	templateName := tplName + tple.templatesSuffix
	if tple.Cached {
		if _, ok := tple.templates[templateName]; ok {
			return tple.templates[templateName]
		}
	}
	
	tplPath := path.Join(tple.templatesDir, tplName+tple.templatesSuffix)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		logger.Warn(err.Error())
	}
	
	if utils.PathExist(path.Join(tple.templatesDir, COMM_TPL_INCLUDES)) {
		if _, err := tpl.ParseGlob(path.Join(tple.templatesDir, COMM_TPL_INCLUDES, "*"+tplName+tple.templatesSuffix)); err != nil {
			logger.Warn(err.Error())
		}
	}
	
	if utils.PathExist(path.Join(tple.templatesDir, COMM_TPL_LAYOUTS)) {
		if _, err := tpl.ParseGlob(path.Join(tple.templatesDir, COMM_TPL_LAYOUTS, "*"+tplName+tple.templatesSuffix)); err != nil {
			logger.Warn(err.Error())
		}
	}
	
	if utils.PathExist(path.Join(tple.templatesDir, COMM_TPL_COMMONS)) {
		if _, err := tpl.ParseGlob(path.Join(tple.templatesDir, COMM_TPL_COMMONS, "*"+tplName+tple.templatesSuffix)); err != nil {
			logger.Warn(err.Error())
		}
	}
	
	logger.Debug("Load Template: %s", templateName)
	if tple.Cached {
		tple.templates[templateName] = tpl
	}
	
	return tpl
}

func (tple *TemplateEngine) Render(tplName string, data map[string]interface{}) (b []byte, err error) {
	tpl := tple.CreateTemplate(tplName)
	
	var buf bytes.Buffer
	e := tpl.Execute(&buf, data)
	//e := tpl.ExecuteTemplate(&buf, tplName+tple.templatesSuffix, data)
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
		tple.Cached = true
		
		layouts, err := filepath.Glob(tple.templatesDir + COMM_TPL_LAYOUTS + "/*" + tple.templatesSuffix)
		if err != nil {
			logger.Error(err.Error())
		} else {
			for _, layout := range layouts {
				tple.templates[filepath.Base(layout)] = template.Must(template.ParseFiles(layout))
				logger.Debug("Load Template Layout: %s, %s", filepath.Base(layout), layout)
			}
		}
		
		includes, err := filepath.Glob(tple.templatesDir + COMM_TPL_INCLUDES + "/*" + tple.templatesSuffix)
		if err != nil {
			logger.Error(err.Error())
		} else {
			for _, include := range includes {
				tple.templates[filepath.Base(include)] = template.Must(template.ParseFiles(include))
				logger.Debug("Load Template Include: %s, %s", filepath.Base(include), include)
			}
		}
		
		commons, err := filepath.Glob(tple.templatesDir + COMM_TPL_COMMONS + "/*" + tple.templatesSuffix)
		if err != nil {
			logger.Error(err.Error())
		} else {
			for _, commons := range commons {
				tple.templates[filepath.Base(commons)] = template.Must(template.ParseFiles(commons))
				logger.Debug("Load Template Commons: %s, %s", filepath.Base(commons), commons)
			}
		}
		
		tple.IsInit = true
	}
}
