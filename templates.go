package getrational

import (
	"encoding/json"
	"html/template"
	"log"
)

func safeHTML(text string) template.HTML {
	return template.HTML(text)
}

func formatJSON(data interface{}) template.JS {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
	}

	return template.JS(bytes)
}

func loadTemplates() (*template.Template, error) {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML": safeHTML,
		"json":     formatJSON,
	})

	templ, err := templ.ParseGlob("templates/*")
	if err != nil {
		return nil, err
	}

	templ = templ.Funcs(template.FuncMap{
		"safeHTML": safeHTML,
	})

	return templ, nil
}
