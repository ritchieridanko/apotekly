package oauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type OAuth struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

// NOTE: Set endpoint to nil to use OIDC provided endpoint
func NewOAuth(ctx context.Context, cfg *configs.OAuth, e *oauth2.Endpoint, l *zap.Logger) (*OAuth, error) {
	p, err := oidc.NewProvider(ctx, cfg.OIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s oauth: %w", cfg.Name, err)
	}

	endpoint := p.Endpoint()
	if e != nil {
		endpoint = *e
	}

	scopes := append(strings.Split(cfg.Scopes, ","), oidc.ScopeOpenID)

	l.Sugar().Infof(
		"[%s-OAUTH] connected (issuer=%s, scopes=[%s])",
		strings.ToUpper(cfg.Name), cfg.OIDCIssuer, cfg.Scopes,
	)
	return &OAuth{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.Secret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       scopes,
			Endpoint:     endpoint,
		},
		verifier: p.Verifier(
			&oidc.Config{
				ClientID:        cfg.ClientID,
				SkipIssuerCheck: cfg.OIDCIssuer == "https://login.microsoftonline.com/common/v2.0",
			},
		),
	}, nil
}

func (o *OAuth) AuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	return o.config.AuthCodeURL(state, opts...)
}

func (o *OAuth) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*Claim, error) {
	token, err := o.config.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token missing from oauth token")
	}

	return o.VerifyIDToken(ctx, rawIDToken)
}

func (o *OAuth) VerifyIDToken(ctx context.Context, rawIDToken string) (*Claim, error) {
	idToken, err := o.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var c Claim
	if err := idToken.Claims(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
