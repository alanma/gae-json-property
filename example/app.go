package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/matthewbelisle-wf/jsonproperty"
	"net/http"
	"time"
)

type Person struct {
	Name     string                    `datastore:"name"`
	Birth    time.Time                 `datastore:"birth"`
	MetaData jsonproperty.JsonProperty `datastore:"-" jsonproperty:"metadata"`
}

func (person *Person) Load(c <-chan datastore.Property) error {
	c, err := jsonproperty.LoadJsonProperties(person, c)
	if err != nil {
		return err
	}
	return datastore.LoadStruct(person, c)
}

func (person *Person) Save(c chan<- datastore.Property) error {
	c, err := jsonproperty.SaveJsonProperties(person, c)
	if err != nil {
		return err
	}
	return datastore.SaveStruct(person, c)
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	person := Person{
		Name:  "Satoshi Nakamoto",
		Birth: time.Date(1215, time.June, 15, 12, 0, 0, 0, time.UTC),
		MetaData: jsonproperty.JsonProperty{
			"favorite_things": []string{
				"life",
				"liberty",
				"pursuit of happiness",
			},
			"point_last_seen": nil,
		},
	}
	c := appengine.NewContext(r)
	key, _ := datastore.Put(c, datastore.NewIncompleteKey(c, "Person", nil), &person)
	person2 := Person{}
	datastore.Get(c, key, &person2)
	bytes, _ := json.Marshal(person2)
	w.Write(bytes)
}
