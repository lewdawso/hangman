package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
	server = "localhost"
	port   = "12345"
)

var help = `
USAGE: hangman COMMAND

COMMAND:
	new                      create a new game
	list    (all|id)         list games / game
	guess   id  [a-z]        guess
	help
`

func usage() {
	fmt.Printf(help)
}

func main() {

	regIndex := regexp.MustCompile("[0-9][0-9]*")
	regChar := regexp.MustCompile("[a-z]")

	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	switch args[1] {
	case "new":
		newGame()
	case "guess":
		if len(args) != 4 {
			usage()
			break
		}
		if !regIndex.MatchString(args[2]) {
			usage()
			break
		}
		if !regChar.MatchString(args[3]) {
			usage()
			break
		}
		guess(args[2], args[3])
	case "list":
		if len(args) != 3 {
			usage()
			break
		}
		if args[2] == "all" {
			listGames()
			break
		}
		if regIndex.MatchString(args[2]) {
			listGame(args[2])
			break
		}
		usage()
	case "help":
		usage()
	default:
		usage()
	}
}
