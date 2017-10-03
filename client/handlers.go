package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var (
	newGameMessage = "Successfully created new game"
	noGamesMessage = "There are no games"
)

func newGame() {

	addr := fmt.Sprintf("http://%s:%s/games", server, port)

	resp, err := netClient.Post(addr, "", strings.NewReader(""))
	if err != nil {
		log.Fatalf("failed to create a new game: %s", err)
	}
	defer resp.Body.Close()

	games, err := processResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(newGameMessage)
	printGames(games)
}

func guess(index, char string) {

	addr := fmt.Sprintf("http://%s:%s/games/%s", server, port, index)

	dat := map[string]string{"guess": char}
	raw, err := json.Marshal(dat)
	if err != nil {
		log.Fatalf("failed to marshal json data: %s", err)
	}

	resp, err := netClient.Post(addr, "application/json", bytes.NewReader(raw))
	if err != nil {
		log.Fatalf("failed to guess character: %s", err)
	}
	defer resp.Body.Close()

	games, err := processResponse(resp)
	if err != nil {
		fmt.Printf("\n%s\n", err)
		return
	}

	printGames(games)
}

func listGames() {

	addr := fmt.Sprintf("http://%s:%s/games", server, port)

	resp, err := netClient.Get(addr)
	if err != nil {
		log.Fatalf("failed to get list of games: %s", err)
	}
	defer resp.Body.Close()

	games, err := processResponse(resp)
	if err != nil {
		fmt.Printf("\n%s", err)
		return
	}

	if len(games) == 0 {
		fmt.Printf("\n%s\n", noGamesMessage)
		return
	}
	printGames(games)
}

func listGame(index string) {

	addr := fmt.Sprintf("http://%s:%s/games/%s", server, port, index)

	resp, err := netClient.Get(addr)
	if err != nil {
		log.Fatalf("failed to get game: %s", err)
	}
	defer resp.Body.Close()

	games, err := processResponse(resp)
	if err != nil {
		fmt.Printf("\n%s", err)
		return
	}
	printGames(games)
}
