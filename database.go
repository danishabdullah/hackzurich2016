package getrational

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

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
	Get(r *http.Request, id string) (GameEntity, error)
	List(r *http.Request, uid string) ([]GameEntity, error)
	Last(r *http.Request, uid string) (*GameEntity, error)
}

type gameDatabase struct{}

type GameEntity struct {
	ID      string    `json:"id"`
	UserID  string    `json:"uid"`
	Time    time.Time `json:"time"`
	Answers []Answer  `json:"answers"`
}

func (db *gameDatabase) Save(r *http.Request, userID, id string, game []Answer) error {
	ctx := appengine.NewContext(r)

	e := &GameEntity{
		ID:      id,
		UserID:  userID,
		Time:    time.Now(),
		Answers: game,
	}

	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	if _, err := datastore.Put(ctx, k, e); err != nil {
		return err
	}

	return nil
}

func (db *gameDatabase) Get(r *http.Request, id string) (GameEntity, error) {
	ctx := appengine.NewContext(r)

	k := datastore.NewKey(ctx, "Game", id, 0, nil)
	var e GameEntity
	if err := datastore.Get(ctx, k, &e); err != nil {
		return GameEntity{}, err
	}

	return e, nil
}

func (db *gameDatabase) List(r *http.Request, uid string) ([]GameEntity, error) {
	ctx := appengine.NewContext(r)

	var result []GameEntity
	q := datastore.NewQuery("Game").Filter("UserID =", uid).Order("ID")
	for t := q.Run(ctx); ; {
		var e GameEntity

		_, err := t.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return []GameEntity{}, err
		}

		result = append(result, e)
	}
	return result, nil
}

func (db *gameDatabase) Last(r *http.Request, uid string) (*GameEntity, error) {
	ctx := appengine.NewContext(r)

	q := datastore.NewQuery("Game").Filter("UserID =", uid).Order("-Time").Limit(1)

	result := new(GameEntity)
	t := q.Run(ctx)

	_, err := t.Next(result)
	if err == datastore.Done {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
