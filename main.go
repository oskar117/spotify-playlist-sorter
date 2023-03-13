// This example demonstrates how to authenticate with Spotify using the authorizati// This example demonstrates how to authenticate with Spotify using the authorization code flow with PKCE.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//     - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
	loc_spotify "github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/zalando/go-keyring"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
		), spotifyauth.WithClientID("0673725a49f845f0b2ee585d87c0df67"))
	ch              = make(chan *spotify.Client, 1)
	state           = "abc123"
	codeVerifier, _ = cv.CreateCodeVerifier()
	codeChallenge   = codeVerifier.CodeChallengeS256()
	service         = "spotify-playlist-sorter"
	user            = "user"
	password        = "secret"
	// These should be randomly generated for each request
	//  More information on generating these can be found here,
	// https://www.oauth.com/playground/authorization-code-with-pkce.html
	// codeVerifier  = "w0HfYrKnG8AihqYHA9_XUPTIcqEXQvCQfOF2IitRgmlF43YWJ8dy2b49ZUwVUOR.YnvzVoTBL57BwIhM4ouSa~tdf0eE_OmiMC_ESCcVOe7maSLIk9IOdBhRstAxjCl7"
	// codeChallenge = "ZhZJzPQXYBMjH8FlGAdYK5AndohLzFfZT-8J7biT7ig"
)

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	tokenString, err := keyring.Get(service, user)
	if err == nil {
		fmt.Println("Loading token from keyring")
		var token oauth2.Token
		json.Unmarshal([]byte(tokenString), &token)
		loadClient(&token)
	} else {
		url := auth.AuthURL(state,
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	}

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	spotifyUser, err := client.CurrentUser(context.Background())
	if err != nil {
		keyring.Delete(service, user)
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", spotifyUser.ID)
	playlistPage, _ := client.GetPlaylistsForUser(context.Background(), spotifyUser.ID)

	var playlistId spotify.ID
	var artists []*sorter.Artist

	for _, playlist := range playlistPage.Playlists {
		if playlist.Owner.ID == spotifyUser.ID && playlist.Name == "asdf" {
			playlistId = playlist.ID
			artists = loc_spotify.FetchArtists(client, playlistId)
		}
	}
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(tui.InitialModel(artists, playlistId, client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	os.Exit(6)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier.String()))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	fmt.Println("Saving token...")
	fmt.Fprintf(w, "Login Completed!")
	tokenAsString, _ := json.Marshal(tok)
	keyring.Set(service, user, string(tokenAsString))
	loadClient(tok)
}

func loadClient(token *oauth2.Token) {
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(context.Background(), token))
	fmt.Println("Login Completed!")
	if m, _ := time.ParseDuration("5m30s"); time.Until(token.Expiry) < m {
		newToken, _ := client.Token()
		tokenAsString, _ := json.Marshal(newToken)
		keyring.Set(service, user, string(tokenAsString))
	}
	ch <- client
	close(ch)
}
