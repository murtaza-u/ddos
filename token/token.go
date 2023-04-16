package token

import (
	"errors"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

var (
	// DefaultExpiry is the default token expiration time if nothing is
	// provided.
	DefaultExpiry = time.Now().Add(time.Minute * 5)

	ErrInvalidNB        = errors.New("expiry cannot be shorter than NB")
	ErrMissingBody      = errors.New("missing body in params")
	ErrMissingSecretKey = errors.New("missing secret key in params")

	protocol = paseto.V4Local
)

// Params contains parameters to construct a paseto token.
type Params struct {
	Expiry    time.Time
	NotBefore time.Time
	Issuer    string
	Audience  string
	Body      map[string]any
}

func (p *Params) validate() error {
	if p.Expiry.IsZero() {
		p.Expiry = DefaultExpiry
	}

	if p.Expiry.Sub(p.NotBefore) < 0 {
		return ErrInvalidNB
	}

	if p.Body == nil || len(p.Body) == 0 {
		return ErrMissingBody
	}

	return nil
}

type token struct {
	paseto.Token
}

// New returns a new paseto token from the provided parameters.
func New(p Params) (*token, error) {
	err := p.validate()
	if err != nil {
		return nil, err
	}

	t := new(token)
	t.Token = paseto.NewToken()
	t.SetExpiration(p.Expiry)
	t.SetIssuedAt(time.Now())
	t.SetNotBefore(p.NotBefore)
	t.SetIssuer(p.Issuer)
	t.SetAudience(p.Audience)

	for k, v := range p.Body {
		err := t.Set(k, v)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Encrypt encrypts and signs the token. The resulting encoded string
// can be given to the client.
func (t *token) Encrypt(secret string, implict []byte) (string, error) {
	key, err := paseto.V4SymmetricKeyFromHex(secret)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create symmetric key from provided secret: %w",
			err,
		)
	}

	return t.V4Encrypt(key, implict), nil
}

// IsExpired checks whether the token has expired. It assumes that the
// "exp" field is a valid RFC3339 compliant time.
func (t *token) IsExpired() bool {
	exp, _ := t.GetExpiration()
	return exp.Before(time.Now())
}

// Refresh refreshes the expired token by updating the expiration to
// (exp - iat)
func (t *token) Refresh() error {
	exp, err := t.GetExpiration()
	if err != nil {
		return fmt.Errorf(
			"failed to get exp from token: %w", err,
		)
	}

	if !t.IsExpired() {
		return fmt.Errorf("token has not yet expired")
	}

	iat, err := t.GetIssuedAt()
	if err != nil {
		return fmt.Errorf(
			"failed to get iat from token: %w", err,
		)
	}

	t.SetIssuedAt(time.Now())
	t.SetExpiration(time.Now().Add(exp.Sub(iat)))
	return nil
}
