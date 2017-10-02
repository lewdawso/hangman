package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	route{
		"listGames",
		"GET",
		"/games",
		listGamesHandler,
	},
	route{
		"listGame",
		"GET",
		"/games/{id:[0-9]*}",
		listGameHandler,
	},
	route{
		"newGame",
		"POST",
		"/games",
		newGameHandler,
	},
	route{
		"guess",
		"POST",
		"/games/{id:[0-9]*}",
		guessHandler,
	},
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
