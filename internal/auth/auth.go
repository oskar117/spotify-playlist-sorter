package auth
//
// import (
// 	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
// 	"github.com/zmb3/spotify"
// 	spotifyauth "github.com/zmb3/spotify/v2/auth"
// 	"golang.org/x/oauth2"
// )
//
// const redirectURI = "http://localhost:8080/callback"
//
// var (
// 	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
// 		spotifyauth.WithScopes(
// 			spotifyauth.ScopeUserReadPrivate,
// 			spotifyauth.ScopePlaylistModifyPrivate,
// 			spotifyauth.ScopePlaylistModifyPublic,
// 		), spotifyauth.WithClientID("0673725a49f845f0b2ee585d87c0df67"))
// 	ch              = make(chan *spotify.Client, 1)
// 	state           = "abc123"
// 	codeVerifier, _ = cv.CreateCodeVerifier()
// 	codeChallenge   = codeVerifier.CodeChallengeS256()
// 	service         = "spotify-playlist-sorter"
// 	user            = "user"
// 	password        = "secret"
// )
//
// func GetOauthToken(resultChannel chan oauth2.Token) oauth2.Token {
// }
