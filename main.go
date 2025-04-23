package main

import (
	"context"
	"fmt"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"log"
)

func main() {
	// Credenciales de OAuth
	apiKey := "JjQ8dRFj0Lmpjkdoz3paYE8ua"
	apiSecret := "vGZIguoxPheDGzENhfvTL05X8LwF54L51obUslSlGWI3d6W41W"
	accessToken := "1914708842019217408-hcbtjAuN9y3JvO7WJwrD3cHHk9kr1p"
	accessSecret := "hUYZRuwcO4ESl4AL7OJImrPDvXVgrYevW4YiLtgxBErwP"

	// Crear cliente OAuth 1.0a
	client, err := newOAuth1Client(apiKey, apiSecret, accessToken, accessSecret)
	if err != nil {
		log.Fatalf("Error al crear el cliente OAuth: %v", err)
	}

	// Texto del tweet
	tweetText := "¡Hola, mundo! Este es un tweet publicado desde Go."

	// Publicar el tweet
	tweetID, err := postTweet(client, tweetText)
	if err != nil {
		log.Fatalf("Error al publicar el tweet: %v", err)
	}

	fmt.Printf("Tweet publicado con éxito. ID: %s\n", tweetID)
}

func newOAuth1Client(apiKey, apiSecret, accessToken, accessSecret string) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		APIKey:               apiKey,
		APIKeySecret:         apiSecret,
		OAuthToken:           accessToken,
		OAuthTokenSecret:     accessSecret,
	}

	return gotwi.NewClient(in)
}

func postTweet(client *gotwi.Client, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), client, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
