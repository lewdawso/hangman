package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lewdawso/hangman/hangman"
	"io/ioutil"
	"log"
	"net/http"
)

func processResponse(resp *http.Response) (map[int]hangman.Game, error) {

	var games map[int]hangman.Game
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read data from response body")
	}

	if resp.StatusCode != http.StatusOK {
		return games, errors.New(string(dat))
	}

	err = json.Unmarshal(dat, &games)
	if err != nil {
		log.Fatalf("failed to unmarshal game data")
	}
	return games, nil
}

func printGames(games map[int]hangman.Game) {
	for i, game := range games {
		fmt.Printf("### Game %d ### \n\n", i)
		fmt.Println(game.String())
	}
}
