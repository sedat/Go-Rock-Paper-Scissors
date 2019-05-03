package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var rounds int
var sessionToken string
var gameRound int
var totalRounds int
var gameScore = make(map[string]int)
var roundResults = []string{}

var types = map[int]string{
	0: "ROCK",
	1: "PAPER",
	2: "SCISSORS",
}

func main() {
	http.HandleFunc("/newGame", newGameHandler)
	http.HandleFunc("/play", playHandler)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func newGameHandler(w http.ResponseWriter, req *http.Request) {
	if rounds, _ = strconv.Atoi(req.URL.Query().Get("rounds")); rounds > 0 {
		rand.Seed(time.Now().UnixNano())
		sessionToken = strconv.Itoa(rand.Intn(100))

		roundResults = make([]string, 0)
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
		w.Write([]byte("Rounds input in wrong format\n"))
	}

}

func playHandler(w http.ResponseWriter, req *http.Request) {
	myChoice := strings.ToUpper(req.FormValue("choose"))
	if getIndex(myChoice) == -1 {
		w.Write([]byte("Choose is in wrong format\n"))
		return
	}

	if req.URL.Query().Get("id") == sessionToken && gameRound < rounds {
		gameRound++
		againstIndex := rand.Intn(3) // Karşı el
		outp := "-> ROUND " + strconv.Itoa(gameRound) + "\n\nme: " + types[againstIndex] + "\nyou: " + myChoice + "\n"
		outp += decideWinner(myChoice, types[againstIndex])
		roundResults = append(roundResults, types[againstIndex]+" vs "+myChoice)
		if gameRound == rounds {
			outp += "\n-> GAME COMPLETED\n"
			for k, v := range roundResults {
				outp += "Round " + strconv.Itoa(k+1) + ": " + v + "\n"
			}
			outp += "\n" + strconv.Itoa(gameScore["computer"]) + " vs " + strconv.Itoa(gameScore["player"])
			if gameScore["player"] > gameScore["computer"] {
				outp += "\nYOU WON\n"
			} else if gameScore["player"] < gameScore["computer"] {
				outp += "\nYOU LOST\n"
			} else {
				outp += "\nIT'S A TIE\n"
			}
		} else {
			outp += "\nThere is " + strconv.Itoa((rounds - gameRound)) + " more round.\n"
		}
		w.Write([]byte(outp))
	} else {
		w.Write([]byte("New game hasn't started or wrong session id\n"))
	}
}

func decideWinner(chosen, against string) string {
	choice := getIndex(chosen)
	versus := getIndex(against)

	if (choice+1)%3 == versus {
		gameScore["computer"]++
		return "\nYOU LOST THIS ROUND\n"
	} else if choice == versus {
		return "\nTIE THIS ROUND!!\n"
	}
	gameScore["player"]++
	return "\nYOU WON THIS ROUND!!\n"

}

func getIndex(word string) int {
	for k, v := range types {
		if word == v {
			return k
		}
	}
	return -1
}
