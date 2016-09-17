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
	if a.Correct() {
		return "success"
	}

	return "danger"
}

func answerEvaluation(answers []Answer) string {
	correct := correctAnswers(answers)
	return evaluateConfidence(int(correct), len(answers), ExpectedConfidence)
}

func correctAnswers(answers []Answer) float64 {
	correct := 0
	for _, a := range answers {
		if a.Correct() {
			correct++
		}
	}
	return float64(correct)
}

func targetScore(answers []Answer) float64 {
	return float64(len(answers)) * ExpectedConfidence
}

func loadTemplates() (*template.Template, error) {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML":   safeHTML,
		"json":       formatJSON,
		"tableClass": tableClass,
		"evaluation": answerEvaluation,
		"correct":    correctAnswers,
		"target":     targetScore,
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
