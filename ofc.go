package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var sessions map[int]*Session = make(map[int]*Session)
var nextSessionId int = 0

type Session struct {
	deck []Card
	curCard int
	curSide int
}

func getDeckJson(addr string) ([]byte, error) {
	response, err := http.Get(addr)
	if err != nil {
		fmt.Println("get_error")
		return []byte(""), err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read_error")
		return []byte(""), err
	}
	return body, nil
}

func (session *Session) loadDeck(deckAddress string) error {
	deckString, err := getDeckJson(deckAddress)
	if err != nil {
		return err
	}
	var orderedDeck []Card
	err = json.Unmarshal(deckString, &orderedDeck)
	if err != nil {
		return err
	}
	deckLength := len(orderedDeck)
	session.deck = make([]Card, deckLength, deckLength)
	order := rand.Perm(deckLength)
	for i, j := range order {
		session.deck[i] = orderedDeck[j]
	}
	return nil
}

func (session *Session) next() string {
	if session.curCard == len(session.deck) {
		return "Done!"
	}
	if session.curSide == 0 {
		text := session.deck[session.curCard].Front
		session.curSide++
		return text
	}
	text := session.deck[session.curCard].Back
	session.curSide = 0
	session.curCard++
	return text
}

type Card struct {
	Front   string `json:"front"`
	Back    string `json:"back"`
	Learned bool   `json:"learned"`
}

type Page struct {
	CardText string
	Link string
}

func handler(w http.ResponseWriter, r *http.Request) {
	firstPath := strings.Split(r.URL.Path, "/")[1]
	sessionId, err := strconv.Atoi(firstPath)

	// We need to create a new session ID
	if err != nil {
		sessionId = nextSessionId
		nextSessionId++
		sessions[sessionId] = &Session{[]Card{}, 0, 0}
		err := (*sessions[sessionId]).loadDeck("http://" + r.URL.Path[1:])
		if err != nil {
			fmt.Println(err)
		}
	}

	session := sessions[sessionId]
	link := fmt.Sprintf("/%d", sessionId)
	page := &Page{CardText: session.next(), Link: link}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, page)
}

func nullHandler(w http.ResponseWriter, r *http.Request) {}

func main() {

	http.HandleFunc("/favicon.ico", nullHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
