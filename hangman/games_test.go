package hangman

//games for testing

//the idea here is to fire valid requests at a game and check that errors are returned when expected and/or the game completes with the expected state

var games = map[string]gameTest{
	"simpleSuccess": gameTest{
		Reply: Reply{
			Error: nil,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 5,
				Found:   6,
				Misses:  []byte("t"),
				Masked:  []byte(testWord),
				State:   success,
			},
		},
		Sequence: "btlargh",
	},
	"allGuessesWrong": gameTest{
		Reply: Reply{
			Error: nil,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 0,
				Found:   0,
				Misses:  []byte("cdefij"),
				Masked:  []byte("______"),
				State:   fail,
			},
		},
		Sequence: "cdefij",
	},
	"alternateCorrectAndWrongSuccess": gameTest{
		Reply: Reply{
			Error: nil,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 1,
				Found:   6,
				Misses:  []byte("cdefi"),
				Masked:  []byte(testWord),
				State:   success,
			},
		},
		Sequence: "bcldaerfgih",
	},
	"alternateCorrectAndWrongFail": gameTest{
		Reply: Reply{
			Error: nil,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 0,
				Found:   5,
				Misses:  []byte("cdefij"),
				Masked:  []byte("blarg_"),
				State:   fail,
			},
		},
		Sequence: "bcldaerfgij",
	},
	"tryAfterFailure": gameTest{
		Reply: Reply{
			Error: errGameCompleted,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 0,
				Found:   3,
				Misses:  []byte("zywqts"),
				Masked:  []byte("bla___"),
				State:   fail,
			},
		},
		Sequence: "blazywqtsf",
	},
	"tryAfterSuccess": gameTest{
		Reply: Reply{
			Error: errGameCompleted,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 4,
				Found:   6,
				Misses:  []byte("cd"),
				Masked:  []byte("blargh"),
				State:   success,
			},
		},
		Sequence: "cdblarghb",
	},
	"repeatSuccessCharacter": gameTest{
		Reply: Reply{
			Error: errAlreadyTriedSuccess,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 2,
				Found:   4,
				Misses:  []byte("zywx"),
				Masked:  []byte("blar__"),
				State:   active,
			},
		},
		Sequence: "blarzywxb",
	},
	"repeatFailedCharacter": gameTest{
		Reply: Reply{
			Error: errAlreadyTriedFail,
			Game: Game{
				word:    emptyWordReply,
				Strikes: 3,
				Found:   3,
				Misses:  []byte("zyx"),
				Masked:  []byte("___rgh"),
				State:   active,
			},
		},
		Sequence: "zyxrghz",
	},
}
