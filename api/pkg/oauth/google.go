package oauth

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

type googleVerifier struct {
	audience string
}

func newGoogleVerifier(audience string) TokenVerifier {
	return &googleVerifier{audience: audience}
}

func (g *googleVerifier) Verify(ctx context.Context, token string) (*Claims, error) {
	payload, err := idtoken.Validate(ctx, token, g.audience)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	if email == "" {
		return nil, fmt.Errorf("token missing email claim")
	}

	return &Claims{
		Subject: payload.Subject,
		Email:   email,
		Name:    name,
	}, nil
}
