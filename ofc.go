package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/tsuru/gnuflag"
)

type Card struct {
	Front   string `json:"front"`
	Back    string `json:"back"`
	Learned bool   `json:"learned"`
}

func main() {

	var reverse bool
	gnuflag.BoolVar(&reverse, "reverse", false, "Allow cards to be presented in reverse.")
	gnuflag.Parse(true)
	var args []string
	if gnuflag.NArg() == 0 {
		args = []string{"example.json"}
	} else {
		args = gnuflag.Args()
	}

	var deck []Card
	deckString, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.Unmarshal(deckString, &deck)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	order := rand.Perm(len(deck))
	for _, i := range order {
		card := deck[i]
		var a string
		var b string
		if reverse {
			if rand.Intn(2) == 0 {
				a = card.Front
				b = card.Back
			} else {
				a = card.Back
				b = card.Front
			}
		} else {
			a = card.Front
			b = card.Back
		}
		fmt.Print(a)
		fmt.Print("\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Print(b)
		fmt.Print("\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
