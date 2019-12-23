package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/just1689/pg-gateway/db"
	"github.com/just1689/pg-gateway/query"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var LB = []byte("[")
var RB = []byte("]")
var COMMA = []byte(",")

func HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, DELETE, POST, OPTIONS")
}

func HandlePatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}/{id}", http.StatusBadRequest)
		return
	}
	field := vars["field"]
	if field == "" {
		http.Error(w, "You need to supply an field", http.StatusBadRequest)
		return
	}
	id := vars["id"]
	if id == "" {
		http.Error(w, "You need to supply an id: /{entity}/{id}", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not read post body", http.StatusBadRequest)
		return
	}

	item := db.Insertable{}
	err = json.Unmarshal(b, &item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not unmarshal item from body", http.StatusBadRequest)
		return
	}

	err = db.Update(entity, field, id, item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not update", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}/{id}", http.StatusBadRequest)
		return
	}
	field := vars["field"]
	if field == "" {
		http.Error(w, "You need to supply an field", http.StatusBadRequest)
		return
	}
	id := vars["id"]
	if id == "" {
		http.Error(w, "You need to supply an id: /{entity}/{id}", http.StatusBadRequest)
		return
	}

	err := db.Delete(entity, field, id)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleGetMany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity", http.StatusBadRequest)
		return
	}

	q, err := query.BuildQueryFromURL(r.URL.String()[1:])
	if err != nil {
		http.Error(w, "Could not build query from http request", http.StatusBadRequest)
		return
	}
	c, err := db.GetByQuery(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "Application/json")
	w.WriteHeader(http.StatusOK)
	rows := 0
	w.Write(LB)
	for row := range c {
		rows++
		if rows > 1 {
			w.Write(COMMA)
		}
		w.Write(row)

	}
	w.Write(RB)

}

func HandleInsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "" {
		http.Error(w, "You need to supply an entity: /{entity}", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not read post body", http.StatusBadRequest)
		return
	}

	item := db.Insertable{}
	err = json.Unmarshal(b, &item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not unmarshal item from body", http.StatusBadRequest)
		return
	}

	err = db.Insert(entity, item)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, "Could not insert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
