package main

import (
	"errors"
	"fmt"
	"github.com/lewdawso/hangman/hangman"
	"log"
	"net/http"
)

const (
	bindAddress = "localhost"
	bindPort    = "12345"
)

var (
	errInvalidGame = errors.New("that game doesn't exist")
)

type req int

const (
	listGame req = iota
	listGames
	newGame
	guess
)

//request and reply structures for interaction between HTTP handlers and server
type request struct {
	kind  req
	value string
	id    int
	done  chan reply
}

type reply struct {
	games map[int]hangman.Game
	err   error
}

type gameChan struct {
	in  chan hangman.Request
	out chan hangman.Reply
}

func (c *gameChan) getState() hangman.Game {
	req := hangman.Request{
		Type:  hangman.State,
		Value: "",
	}
	c.in <- req
	rep := <-c.out

	return rep.Game
}

func (c *gameChan) guess(ch string) (hangman.Game, error) {
	req := hangman.Request{
		Type:  hangman.Guess,
		Value: ch,
	}
	c.in <- req
	rep := <-c.out

	return rep.Game, rep.Error
}

var handlerChan chan request

func main() {

	var gameCount int

	//this is all we need to store - an index + a game's input and output channels
	games := make(map[int]gameChan)

	//HTTP handlers send requests to this channel
	handlerChan = make(chan request)

	/*
		a user can do three things:

		1. start a new game
		2. view a list of all games (active and completed)
		3. interact with any game (again, both active and completed)

		Each of these operations will effect the outcome of another => only one can happen at a time

		Concurrency is handled at this layer i.e. by serialising requests from HTTP handlers

	*/

	//this goroutine services requests received on handlerChan
	go func() {
		for {
			r := <-handlerChan

			switch r.kind {
			case newGame:
				index := gameCount
				rep := reply{}
				in, out := hangman.NewGame()
				games[gameCount] = gameChan{
					in,
					out,
				}
				gameCount++
				//now get the state
				chans, _ := games[index]
				game := chans.getState()
				rep.games = map[int]hangman.Game{
					index: game,
				}
				r.done <- rep
			case listGames:
				rep := reply{}
				list := make(map[int]hangman.Game)
				for i, chans := range games {
					list[i] = chans.getState()
				}
				rep.games = list
				r.done <- rep
			case listGame:
				id := r.id
				rep := reply{}
				chans, ok := games[id]
				if !ok {
					rep.err = errInvalidGame
					r.done <- rep
					break
				}
				game := chans.getState()
				rep.games = map[int]hangman.Game{
					id: game,
				}
				r.done <- rep
			case guess:
				id := r.id
				char := r.value
				rep := reply{}
				chans, ok := games[id]
				if !ok {
					rep.err = errInvalidGame
					r.done <- rep
					break
				}
				game, err := chans.guess(char)
				rep.err = err
				rep.games = map[int]hangman.Game{
					id: game,
				}
				r.done <- rep
			default:
				break
			}
		}
	}()

	router := newRouter()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", bindAddress, bindPort), router))
}
