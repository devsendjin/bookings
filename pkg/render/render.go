package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/devsendjin/bookings/pkg/config"
	"github.com/devsendjin/bookings/pkg/models"
	// "golang.org/x/exp/maps"
)

var templateFunctions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	currentTemplate, ok := templateCache[tmpl]
	// fmt.Println("RenderTemplate - currentTemplate:", currentTemplate)
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	currentTemplate.Execute(buffer, templateData)

	_, err := buffer.WriteTo(w)
	if err != nil {
		fmt.Println("RenderTemplate - Error writing template to browser", err)
	}

	// fmt.Printf("\n\n")
	// fmt.Printf("----------------------------")
	// fmt.Printf("\n\n")
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		// fmt.Println("CreateTemplateCache - page:", page)

		name := filepath.Base(page)

		// fmt.Println("CreateTemplateCache - name:", name)

		templateSet, err := template.New(name).Funcs(templateFunctions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// fmt.Println("CreateTemplateCache - templateSet:", templateSet)

		layoutGlob := "./templates/*.layout.tmpl"

		matches, err := filepath.Glob(layoutGlob)
		if err != nil {
			return templateCache, err
		}

		// fmt.Println("CreateTemplateCache - matches:", matches)

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(layoutGlob)
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	// fmt.Println("CreateTemplateCache - myCache:", templateCache)

	// fmt.Println("pages", pages)

	return templateCache, nil
}
