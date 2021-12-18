package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientId     = "myclient"
	clientSecret = "p78wSx6zRhuAGd7exuuqtV0MVyx1d3TG"
	issuer       = "http://localhost:8080/auth/realms/myrealm"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(rw, "Invalid state", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(rw, "Could not exchange the code by an access token", http.StatusInternalServerError)
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(rw, "Failed to generate the id token", http.StatusInternalServerError)
			return
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(rw, "Could not get user info", http.StatusInternalServerError)
			return
		}

		res := struct {
			AccessToken *oauth2.Token
			IDToken     string
			UserInfo    *oidc.UserInfo
		}{
			AccessToken: token,
			IDToken:     idToken,
			UserInfo:    userInfo,
		}

		data, err := json.Marshal(res)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
