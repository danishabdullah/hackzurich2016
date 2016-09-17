package main

import "html/template"

func safeHTML(text string) template.HTML {
	return template.HTML(text)
}

func loadTemplates() (*template.Template, error) {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML": safeHTML,
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
