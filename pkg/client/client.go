package client

import (
	"context"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func NewSpotifyClient(clientID, clientSecret string) *spotify.Client {
	cfg := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	client := cfg.Client(context.Background())
	newClient := spotify.NewClient(client)
	return &newClient
}
