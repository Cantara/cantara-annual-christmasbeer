package oidc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

type claims interface {
	jwt.StandardClaims
	Valid(*jwt.ValidationHelper) error
}

type Validator struct {
	ar       *jwk.AutoRefresh
	url      string
	issuer   string
	audience string
	ctx      context.Context
}

func InitValidator(url, issuer, aud string, ctx context.Context) (v Validator, err error) {

	//const googleCerts = `https://www.googleapis.com/oauth2/v3/certs`
	//const googleCerts = `https://id.twitch.tv/oauth2/keys`
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(url, jwk.WithMinRefreshInterval(15*time.Minute))
	_, err = ar.Refresh(ctx, url)
	if err != nil {
		err = fmt.Errorf("failed to refresh %s JWKS: %s\n", url, err)
		return
	}
	v = Validator{
		ar:       ar,
		url:      url,
		issuer:   issuer,
		audience: aud,
		ctx:      ctx,
	}
	return
}

type EmailClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.StandardClaims
}

func (v Validator) Validate(tokenString string, claims jwt.Claims) (token *jwt.Token, err error) {
	keyset, err := v.ar.Fetch(v.ctx, v.url)
	if err != nil {
		err = fmt.Errorf("failed to fetch %s JWKS: %s\n", v.url, err)
		return
	}

	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		key, ok := keyset.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}
		var raw interface{}
		return raw, key.Raw(&raw)
	}, jwt.WithIssuer(v.issuer), jwt.WithAudience(v.audience))
	if err != nil {
		return
	}
	return
}
