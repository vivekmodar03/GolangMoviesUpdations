package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

// Init initializes the Firebase Admin SDK
func Init() {
	opt := option.WithCredentialsFile("firebase-credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting firebase auth client: %v\n", err)
	}

	log.Println("Firebase Admin SDK Initialized Successfully")
}