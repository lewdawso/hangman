# hangman

Server implementation of a hangman game.

## Building

To build the server:

```
cd server && go build
```

To build the hangman client:

```
cd client && ./build.sh
```

## Client Usage

```
USAGE: hangman COMMAND

COMMAND:
	new                      create a new game
	list    (all|id)         list games / game
	guess   id  [a-z]        guess
	help
```


## Testing

Only unit tests that test the hangman implementation have been written. To run these do:

```
cd hangman && go test
cd hangman && go test --race
```

The focus here is on testing the interface provided by the "hangman" package.

go test --cover = 89.7%

## Concurrency

Each hangman game lives inside a goroutine. The only way to talk to it is by sending requests to an input channel and
waiting for a response from an output channel. This serialises access to each game.

Similarly, the server needs to be able to handle concurrent HTTP requests. You should not, for example,
be able to guess a character for a game whilst a list of games is being generated - this would lead to inconsistent
results. To deal with this, the "server" also lives inside a goroutine and is exposed only through a channel which the
HTTP handlers can send their requests to.

## To-do

* Integration tests (run server in container, send API requests, compare returned data to expected)
* Server unit tests
