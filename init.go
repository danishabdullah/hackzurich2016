package getrational

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	templ, err := loadTemplates()
	if err != nil {
		log.Fatalf("Can not load templates: %s", err)
	}

	questions, err := readDatabase()
	if err != nil {
		log.Fatalf("Can not read database: %s", err)
	}

	games := make(GameDatabase)

	initHandlers(http.DefaultServeMux, templ, questions, games)
}
