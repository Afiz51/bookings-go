package render

import (
	"bytes"

	"html/template"
	"log"
	"net/http"
	"path/filepath"

	models "github.com/Afiz51/bookings-go/pkg/Models"
	"github.com/Afiz51/bookings-go/pkg/config"
)

var tc = make(map[string]*template.Template)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// renders templaate using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//create a template cache

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all files named *.pagetmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	log.Println("all files named *.page.tmpl: ", pages)

	if err != nil {
		return myCache, err
	}

	//range through all files ending with *page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}
