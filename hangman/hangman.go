package hangman

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
)

//hangman game implementation

const (
	maxTries = 6
)

var (
	emptyWordReply       = []byte("")
	gameWords      words = allWords
)

var (
	errInvalidCharacter    = errors.New("that's not a valid character")
	errAlreadyTriedFail    = errors.New("you've already guessed that letter")
	errAlreadyTriedSuccess = errors.New("you've already found that letter")
	errGameCompleted       = errors.New("the game has already finished")
	errInvalidReqType      = errors.New("invalid request type")
)

//Game is the hangman game structure. "word" is private so that it's only seen by functions internal to the hangman package.
type Game struct {
	word    []byte //server should never see this
	Strikes int    `json:"strikes"`
	Found   int    `json:"found"`
	Misses  []byte `json:"misses"`
	Masked  []byte `json:"masked"`
	State   state  `json:"state"`
}

type state int

const (
	active state = iota
	fail
	success
)

var stateString = map[state]string{
	active:  "IN PROGRESS",
	fail:    "GAME OVER",
	success: "YOU WON",
}

type req int

//You can interact with a game of hangman in two ways:
//1. Guess a letter
//2. Request the game's State.
const (
	Guess req = iota
	State
)

//Request encapsulates both the type of request (guess/state) and any data provided with it i.e. the character to guess.
type Request struct {
	Type  req
	Value string
}

//Reply is the structure returned by a game following a request. It contains a copy of the game structure at that moment plus any error encountered during the processing of the request.
type Reply struct {
	Game  Game
	Error error
}

type words []string

func (w *words) random() []byte {
	index := rand.Intn(len(*w))
	return []byte((*w)[index])
}

//NewGame initialises a new game, launches a goroutine for handling requests and returns an input and output channel for interacting with the game.
func NewGame() (chan Request, chan Reply) {
	g := Game{}
	g.word = gameWords.random()
	g.Strikes = maxTries
	g.Found = 0
	g.State = active

	for i := 0; i < len(g.word); i++ {
		g.Masked = append(g.Masked, "_"[0])
	}

	in := make(chan Request)
	out := make(chan Reply)

	go g.handleRequest(in, out)
	return in, out
}

//for nice print
func (g *Game) String() string {
	ascii := fmt.Sprintf("%s", hangpics[6-g.Strikes])
	masked := fmt.Sprintf("%s\n\n", g.Masked)
	state := fmt.Sprintf("State: %s", stateString[g.State])
	misses := fmt.Sprintf("Incorrect Guesses: %s\n", g.Misses)

	return fmt.Sprintf("%s\t%s%s%s", ascii, masked, misses, state)
}

//business logic of game
func (g *Game) guess(ch string) error {

	correct := false

	if g.State == fail || g.State == success {
		return errGameCompleted
	}

	if len(ch) != 1 {
		return errInvalidCharacter
	}

	//should really handle caps here
	reg := regexp.MustCompile("[a-z]")

	if match := reg.MatchString(ch); match != true {
		return errInvalidCharacter
	}

	for _, v := range g.Misses {
		if string(v) == ch {
			return errAlreadyTriedFail
		}
	}

	for _, v := range g.Masked {
		if string(v) == ch {
			return errAlreadyTriedSuccess
		}
	}

	for i, v := range g.word {
		if string(v) == ch {
			//if it can't happen, don't let it happen
			if g.Masked[i] != "_"[0] {
				log.Fatalf("expected _ but found %s", string(g.Masked[i]))
			}
			g.Masked[i] = ch[0]
			g.Found++
			correct = true
		}
	}

	if !correct && g.Strikes > 0 {
		g.Strikes--
		g.Misses = append(g.Misses, ch[0])
	}

	if len(g.word) == g.Found && g.Strikes >= 0 {
		g.State = success
		return nil
	}

	if g.Strikes == 0 && g.State != success {
		g.State = fail
		return nil
	}

	return nil
}

//The only way to interact with a game is to send "Request" messages to its request handler via an input channel that's returned when the game is initialised.
//The idea is that a game lives in its own goroutine and we only have access to a channel
//to interact with it. The goroutine lives as long as the server runs so that we can continue to query the state.
func (g *Game) handleRequest(in chan Request, out chan Reply) {
	for {
		r := <-in
		var rep Reply
		switch r.Type {
		case Guess:
			err := g.guess(r.Value)

			rep = Reply{
				Error: err,
				Game: Game{
					word:    emptyWordReply,
					Strikes: g.Strikes,
					Found:   g.Found,
					Misses:  g.Misses,
					Masked:  g.Masked,
					State:   g.State,
				},
			}
		case State:
			rep = Reply{
				Error: nil,
				Game: Game{
					word:    emptyWordReply,
					Strikes: g.Strikes,
					Found:   g.Found,
					Misses:  g.Misses,
					Masked:  g.Masked,
					State:   g.State,
				},
			}
		default:
			//invalid req type
			rep = Reply{
				Error: errInvalidReqType,
			}
		}
		out <- rep
	}
}
