package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kousuke1201abe/go-haiku/haiku"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var token = os.Getenv("TOKEN")
var vtoken = os.Getenv("VTOKEN")
var botname = os.Getenv("BOTNAME")
var api = slack.New(token)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: vtoken}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				if ev.User != botname {
					words := haiku.NewWords(ev.Text)
					if words.CheckHaiku() == true {
						api.PostMessage(ev.Channel, slack.MsgOptionText("5 7 5", false))
					}
				}
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":"+port, nil)
}
