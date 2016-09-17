package getrational

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

// Question is the basic data entity.
type Question struct {
	Text      string `json:"text"`
	Unit      string `json:"unit"`
	BoundLow  uint64 `json:"boundLow"`
	BoundHigh uint64 `json:"boundHigh"`
}

func convertRecord(rec []string) (Question, error) {
	if len(rec) != 4 {
		return Question{}, fmt.Errorf("invalid length: %d", len(rec))
	}

	low, err := strconv.ParseFloat(rec[1], 64)
	if err != nil {
		return Question{}, err
	}

	high, err := strconv.ParseFloat(rec[2], 64)
	if err != nil {
		return Question{}, err
	}

	return Question{
		Text:      rec[0],
		Unit:      rec[3],
		BoundLow:  uint64(low),
		BoundHigh: uint64(high),
	}, nil
}

func readDatabase() ([]Question, error) {
	f, err := os.Open("Questions.csv")
	if err != nil {
		return []Question{}, err
	}

	reader := csv.NewReader(f)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return []Question{}, err
	}

	var result []Question
	for _, rec := range records {
		q, err := convertRecord(rec)
		if err != nil {
			log.Printf("Invalid record: %s", err)
			continue
		}

		result = append(result, q)
	}
	return result, nil
}

// QuestionDatabase is the interface for the database containing the questions.
type QuestionDatabase []Question

// SelectRandom selects `num` questions at random from the database.
func (db QuestionDatabase) SelectRandom(num int) []Question {
	if len(db) < num {
		return db
	}

	idx := rand.Perm(len(db))
	var result []Question
	for c, i := range idx {
		if c >= num {
			break
		}
		result = append(result, db[i])
	}

	return result
}
