package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

var rounds string
var sessionToken string
var gameRound int
var totalRounds int
var gameScore = make(map[string]int)
var roundResults = make(map[int]string)

var types = map[int]string{
	0: "rock",
	1: "paper",
	2: "scissors",
}

func main() {
	http.HandleFunc("/newGame", newGameHandler)
	http.HandleFunc("/play", playHandler)
	http.ListenAndServe(":8080", nil)
}

func newGameHandler(w http.ResponseWriter, req *http.Request) {
	if rounds = req.URL.Query().Get("rounds"); rounds != "" {
		totalRounds, _ = strconv.Atoi(rounds)
		sessionToken = strconv.Itoa(rand.Intn(100))

		roundResults = make(map[int]string)
		gameRound = 0
		gameScore["player"], gameScore["computer"] = 0, 0

		// Session ID için cookie kullanılabilir
		// http.SetCookie(w, &http.Cookie{
		// 	Name:  "sessionID",
		// 	Value: sessionToken,
		// })

		w.Write([]byte("New Rock-Paper-Scissors game started\nSession ID = " + sessionToken + "\n\n" +
			"to play Rock: http://localhost:8080/play?choose=rock&id=" + sessionToken +
			"\nto play Paper: http://localhost:8080/play?choose=paper&id=" + sessionToken +
			"\nto play Scissors: http://localhost:8080/play?choose=scissors&id=" + sessionToken + "\n"))

	} else {
		w.Write([]byte("Rounds cannot be empty!"))
	}

}

func playHandler(w http.ResponseWriter, req *http.Request) {

	if req.URL.Query().Get("id") == sessionToken && gameRound < totalRounds {
		againstIndex := rand.Intn(3)
		gameRound++
		outp := "-> ROUND " + strconv.Itoa(gameRound) + "\n\nme: " + req.FormValue("choose") + "\nyou: " + types[againstIndex] + "\n"
		outp += decideWinner(req.FormValue("choose"), types[againstIndex])
		roundResults[gameRound] = req.FormValue("choose") + " vs " + types[againstIndex]
		if gameRound == totalRounds {
			outp += "\n-> GAME COMPLETED\n"
			for k, v := range roundResults {
				outp += "Round " + strconv.Itoa(k) + ": " + v + "\n"
			}
			outp += strconv.Itoa(gameScore["computer"]) + " vs " + strconv.Itoa(gameScore["player"])
			if gameScore["player"] > gameScore["computer"] {
				outp += "\nYOU WON\n"
			}
			if gameScore["player"] < gameScore["computer"] {
				outp += "\nYOU LOST\n"
			} else {
				outp += "\nIT'S A TIE\n"
			}
		} else {
			outp += "\nThere is " + strconv.Itoa((totalRounds - gameRound)) + " more round\n"
		}
		w.Write([]byte(outp))
	} else {
		w.Write([]byte("New game hasn't started or wrong session id"))
	}
}

func decideWinner(chosen string, against string) string {
	choice := getIndex(chosen)
	versus := getIndex(against)

	if (choice+1)%3 == versus {
		gameScore["computer"]++
		return "\nYOU LOST THIS ROUND\n"
	}
	if choice == versus {
		return "\nTIE THIS ROUND!!\n"
	} else {
		gameScore["player"]++
		return "\nYOU WON THIS ROUND!!\n"
	}
}

func getIndex(word string) int {
	for k, v := range types {
		if word == v {
			return k
		}
	}
	return -1
}
