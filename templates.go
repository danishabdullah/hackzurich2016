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

func tableClass(a Answer) string {
	qLow := a.Question.BoundLow
	qHigh := a.Question.BoundHigh
	aLow := a.LowerBound
	aHigh := a.UpperBound
	success := (aLow >= qLow && aHigh <= qHigh) ||
		(aLow <= qLow && aHigh >= qLow) ||
		(aLow <= qHigh && aHigh >= qHigh)

	if success {
		return "success"
	}

	return "danger"
}

func loadTemplates() (*template.Template, error) {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML":   safeHTML,
		"json":       formatJSON,
		"tableClass": tableClass,
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
