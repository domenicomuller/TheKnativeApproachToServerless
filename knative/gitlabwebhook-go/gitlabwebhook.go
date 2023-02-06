package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/webhooks/v6/gitlab"
)

const (
	// Personal Access Token created in github that allows us to make
	// calls into github.
	webhookSecretKey = "WEBHOOK_SECRET"
	path = "/"
)

func main() {
	log.Print("gitwebhook sample started.")
	secretToken := os.Getenv(webhookSecretKey)
	hook, _ := gitlab.New(gitlab.Options.Secret(secretToken))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		log.Print("Handling Pull Request")

		payload, err := hook.Parse(r, gitlab.PushEvents)
		if err != nil {
			
			if err == gitlab.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case gitlab.PushEventPayload:
			pullRequest := payload.(gitlab.PushEventPayload)
			// Do whatever you want from here...
			log.Print("Here it is the request: ")
			fmt.Printf("%+v", pullRequest)
		}
	})
	http.ListenAndServe(":8080", nil)
}