package predictiongame

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
)

func safeHTML(text string) template.HTML {
	return template.HTML(text)
}

func rangeStr(lower, upper float64) string {
	if lower == upper {
		return fmt.Sprintf("%.0f", lower)
	}

	return fmt.Sprintf("%.0f-%.0f", lower, upper)
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

func correctAnswersPercent(answers []Answer) string {
	correct := correctAnswers(answers)
	return fmt.Sprintf("%.0f%%", correct/float64(len(answers))*100)
}

func targetScore(answers []Answer) float64 {
	return float64(len(answers)) * ExpectedConfidence
}

func countHistory(games []GameEntity) (float64, int) {
	correct := 0.0
	count := 0
	for _, g := range games {
		correct += correctAnswers(g.Answers)
		count += len(g.Answers)
	}
	return float64(correct), count
}

func correctAnswersHistory(games []GameEntity) float64 {
	correct, _ := countHistory(games)
	return correct
}

func correctAnswersHistoryPercent(games []GameEntity) string {
	correct, count := countHistory(games)
	return fmt.Sprintf("%.0f%%", correct/float64(count)*100)
}

func targetScoreHistory(games []GameEntity) float64 {
	_, count := countHistory(games)
	return float64(count) * ExpectedConfidence
}

func offset(value, offset int) int {
	return value + offset
}

func loadTemplates() (*template.Template, error) {
	templ := template.New("root").Funcs(template.FuncMap{
		"safeHTML":              safeHTML,
		"rangeStr":              rangeStr,
		"json":                  formatJSON,
		"tableClass":            tableClass,
		"evaluation":            answerEvaluation,
		"correct":               correctAnswers,
		"correctPercent":        correctAnswersPercent,
		"target":                targetScore,
		"correctHistory":        correctAnswersHistory,
		"correctHistoryPercent": correctAnswersHistoryPercent,
		"targetHistory":         targetScoreHistory,
		"offset":                offset,
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
