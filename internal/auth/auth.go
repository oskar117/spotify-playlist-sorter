package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/zalando/go-keyring"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
		), spotifyauth.WithClientID("0673725a49f845f0b2ee585d87c0df67"))
	tokenChannel    = make(chan *oauth2.Token, 1)
	state           = "abc123"
	codeVerifier, _ = cv.CreateCodeVerifier()
	codeChallenge   = codeVerifier.CodeChallengeS256()
	service         = "spotify-playlist-sorter"
	user            = "user"
	password        = "secret"
)

func GetOauthToken() *oauth2.Token {
	startAuthServer()
	tokenString, err := keyring.Get(service, user)
	if err == nil {
		log.Println("Loading token from keyring")
		var token oauth2.Token
		json.Unmarshal([]byte(tokenString), &token)
		return &token
	}
	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	return <-tokenChannel
}

func GetHttpClient() *http.Client {
	return auth.Client(context.Background(), GetOauthToken())
}

func RemoveTokenFromKeyring() {
	log.Println("Removing tokens")
	keyring.Delete(service, user)
}

func UpdateToken(token *oauth2.Token) {
	log.Println("Saving token...")
	tokenAsString, _ := json.Marshal(token)
	keyring.Set(service, user, string(tokenAsString))
}

func startAuthServer() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier.String()))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, state)
	}
	fmt.Fprintf(w, "Login Completed!")
	log.Println("Login completed, returning token!")
	tokenChannel <- tok
	close(tokenChannel)
}

