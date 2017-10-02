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

## Usage

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

## To-do

* Integration tests (run server in container, send API requests, compare returned data to expected)
