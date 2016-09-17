package getrational

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Question is the basic data entity.
type Question struct {
	Text      string  `json:"text"`
	Unit      string  `json:"unit"`
	BoundLow  float64 `json:"boundLow"`
	BoundHigh float64 `json:"boundHigh"`
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
		BoundLow:  low,
		BoundHigh: high,
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

type GameDatabase interface {
	Save(r *http.Request, userID, id string, game []Answer) error
	Get(r *http.Request, id string) ([]Answer, error)
}

type gameDatabase struct{}

type GameEntity struct {
	ID      string
	UserID  string
	Answers []Answer
}

func (db *gameDatabase) Save(r *http.Request, userID, id string, game []Answer) error {
	ctx := appengine.NewContext(r)

	e := &GameEntity{
		ID:      id,
		UserID:  userID,
		Answers: game,
	}

	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	v, err := datastore.Put(ctx, k, e)
	if err != nil {
		return err
	}

	log.Printf("value: %+v", v)
	return nil
}

func (db *gameDatabase) Get(r *http.Request, id string) ([]Answer, error) {
	ctx := appengine.NewContext(r)

	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	e := new(GameEntity)
	if err := datastore.Get(ctx, k, e); err != nil {
		return []Answer{}, err
	}

	return e.Answers, nil
}
