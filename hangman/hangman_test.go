package hangman

import (
	"reflect"
	"testing"
)

var (
	testWord = "blargh"
)

//test games using struct containing the expected response + sequence of attempted characters
type gameTest struct {
	Reply    Reply
	Sequence string
}

//expected initialised state
var expectInit = Reply{
	Error: nil,
	Game: Game{
		word:    emptyWordReply,
		Strikes: maxTries,
		Found:   0,
		Misses:  nil,
		Masked:  []byte("______"),
		State:   active,
	},
}

//does the game structure initialise with what we expect?
func TestInit(t *testing.T) {

	in, out := NewGame()

	req := Request{
		Type:  State,
		Value: "",
	}

	in <- req

	reply := <-out

	if !reflect.DeepEqual(reply, expectInit) {
		t.Errorf("failed to initialise game structure correctly")
	}
}

//try some invalid characters
func TestInvalidChar(t *testing.T) {

	invalidChars := []string{"_", ".", "*", "1", "9", "%", "\n", "0x10", "A", "Z", "", " "}
	for _, v := range invalidChars {
		in, out := NewGame()

		req := Request{
			Type:  Guess,
			Value: v,
		}

		in <- req

		reply := <-out

		if reply.Error != errInvalidCharacter {
			t.Errorf("expected invalid character error for char: %s", v)
		}
	}
}

func TestInvalidReqType(t *testing.T) {

	in, out := NewGame()

	req := Request{
		Type:  2,
		Value: "",
	}

	in <- req

	reply := <-out

	if reply.Error != errInvalidReqType {
		t.Errorf("expected error %s but got %s", errInvalidReqType, reply.Error)
	}
}

func TestWordWithTwoCharSame(t *testing.T) {

	testWord := "aardvark"

	gameWords = []string{testWord}

	game := gameTest{
		Reply: Reply{
			Error: nil,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 5,
				Found:   8,
				Misses:  []byte("z"),
				Masked:  []byte(testWord),
				State:   success,
			},
		},
		Sequence: "azrdvrk",
	}

	runGame("testWordWithTwoCharSame", game, t)
}

func TestGame(t *testing.T) {

	//set gameWords to our testWord for predictable output
	gameWords = []string{testWord}

	for name, game := range games {
		runGame(name, game, t)
	}
}

func runGame(name string, test gameTest, t *testing.T) {
	in, out := NewGame()
	var res Reply
	for _, v := range test.Sequence {
		//keep firing request until we're done
		//don't worry about output until the end
		req := Request{
			Type:  Guess,
			Value: string(v),
		}

		in <- req
		res = <-out
	}

	//now consider the response
	if test.Reply.Error != res.Error {
		t.Errorf("expected error %s but got %s for game %s", test.Reply.Error, res.Error, name)
	}
	if !reflect.DeepEqual(res.Game, test.Reply.Game) {
		t.Errorf("game state differs from expected: %s", name)
	}
}
