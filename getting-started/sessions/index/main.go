// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START getting_started_background_app_main]

// Command index is an HTTP app that displays the number of accesses and a
// translated greeting
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	"github.com/dgrijalva/jwt-go"
)

var greetings = []string{"Hello World", "Hallo Welt", "Hola mundo", "Salut le Monde", "Ciao Mondo"}

// An app holds the Firestore client
type app struct {
	firestoreApp *firebase.App
	authClient   *auth.Client
	ctx          context.Context
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT must be set")
	}
	a, err := newApp(projectID, "index")
	if err != nil {
		log.Fatalf("newApp: %v", err)
	}

	http.HandleFunc("/", a.index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on localhost:%v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// newApp creates a new app.
func newApp(projectID, templateDir string) (*app, error) {
	ctx := context.Background()
	config := firebase.Config{
		DatabaseURL:   "https://sample-248520.firebaseio.com",
		ProjectID:     "sample-248520",
		StorageBucket: "sample-248520.appspot.com",
	}
	firestoreApp, err := firebase.NewApp(ctx, &config)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %v", err)
	}
	authClient, err := firestoreApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}

	return &app{
		firestoreApp: firestoreApp,
		authClient:   authClient,
		ctx:          ctx,
	}, nil
}

// index gives a translated greeting.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	a.getCookie(w, r)
	n := rand.Intn(len(greetings))
	fmt.Fprintf(w, "%d views for %s", 0, greetings[n])
}

func (a *app) getCookie(w http.ResponseWriter, r *http.Request) {
	// Get the ID token sent by the client
	defer r.Body.Close()
	idToken := jwt.New(jwt.SigningMethodHS256).Raw
	// Set session expiration to 5 days.
	expiresIn := time.Hour * 24 * 5
	log.Printf("token %v", idToken)

	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	// To only allow session cookie setting on recent sign-in, auth_time in ID token
	// can be checked to ensure user was recently signed in before creating a session cookie.
	cookie, err := a.authClient.SessionCookie(r.Context(), idToken, expiresIn)
	if err != nil {
		http.Error(w, "Failed to create a session cookie", http.StatusInternalServerError)
		log.Printf("SessionCookie: %v", err)
		return
	}

	// Set cookie policy for session cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    cookie,
		MaxAge:   int(expiresIn.Seconds()),
		HttpOnly: true,
		Secure:   true,
	})
	w.Write([]byte(`{"status": "success"}`))
}

func getIDTokenFromBody(r *http.Request) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var parsedBody struct {
		IDToken string `json:"idToken"`
	}
	err = json.Unmarshal(b, &parsedBody)
	return parsedBody.IDToken, err
}
