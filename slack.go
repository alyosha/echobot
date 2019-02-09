package main

import (
	"context"
	"errors"

	"github.com/nlopes/slack"
)

type signingSecretKey struct{}
type slackClientKey struct{}
type channelKey struct{}

// getSigningSecret is the method used to extract the signing secret from the request context
func getSigningSecret(ctx context.Context) (string, error) {
	val := ctx.Value(signingSecretKey{})
	secret, ok := val.(string)
	if !ok {
		return "", errors.New("error extracting the signing secret from context")
	}

	return secret, nil
}

// withContext embeds values into to the request context
func withContext(ctx context.Context, signingSecret string, client *slack.Client) context.Context {
	return addClient(addSigningSecret(ctx, signingSecret), client)
}

func addSigningSecret(ctx context.Context, signingSecret string) context.Context {
	return context.WithValue(ctx, signingSecretKey{}, signingSecret)
}

func addClient(ctx context.Context, client *slack.Client) context.Context {
	return context.WithValue(ctx, slackClientKey{}, client)
}
