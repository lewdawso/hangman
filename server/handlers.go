package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func newGameHandler(w http.ResponseWriter, r *http.Request) {

	req := request{
		kind:  newGame,
		value: "",
		done:  make(chan reply),
	}

	handlerChan <- req
	rep := <-req.done

	if !checkError(w, rep.err) {
		return
	}

	if err := json.NewEncoder(w).Encode(rep.games); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Fatalf("error marhsalling game: %s", err)
	}
}

func listGamesHandler(w http.ResponseWriter, r *http.Request) {

	req := request{
		kind:  listGames,
		value: "",
		done:  make(chan reply),
	}

	handlerChan <- req
	rep := <-req.done

	if !checkError(w, rep.err) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(rep.games); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Fatalf("error marhsalling list of games: %s", err)
	}
}

func listGameHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var id string
	var ok bool
	if id, ok = vars["id"]; !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	iden, _ := strconv.Atoi(id)

	req := request{
		kind: listGame,
		id:   iden,
		done: make(chan reply),
	}

	handlerChan <- req
	rep := <-req.done

	if !checkError(w, rep.err) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(rep.games); err != nil {
		log.Printf("error marhsalling game: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	dat, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("failed to read data from response body: %s", err)
	}

	var expect map[string]string

	err = json.Unmarshal(dat, &expect)
	if err != nil {
		log.Fatalf("failed to unmarshal json data: %s", err)
	}

	var char string
	var ok bool
	if char, ok = expect["guess"]; !ok {
		log.Printf("received POST at /game/%s not containing data", id)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	iden, _ := strconv.Atoi(id)

	req := request{
		kind:  guess,
		id:    iden,
		value: char,
		done:  make(chan reply),
	}

	handlerChan <- req
	rep := <-req.done

	if !checkError(w, rep.err) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(rep.games); err != nil {
		log.Printf("error marhsalling game: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
	}

}
