package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// TokenSource implements oauth2.TokenSource.
type TokenSource struct {
	// ClusterID represents the cluster ID.
	ClusterID string
	// TokenGenerator is used to generate the token.
	TokenGenerator TokenGenerator
	// Leeway allows refreshing the token before its expiry.
	Leeway time.Duration
}

// Token returns the token.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token, err := t.TokenGenerator.GetWithSTS(context.Background(), t.ClusterID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	return &oauth2.Token{
		AccessToken: token.Token,
		Expiry:      token.Expiration.Add(-t.Leeway),
	}, nil
}
