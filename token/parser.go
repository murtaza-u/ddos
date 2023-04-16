package token

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

// Decrypt decrypts the encrypted token
func Decrypt(enc, secret string, implicit []byte) (*token, error) {
	key, err := paseto.V4SymmetricKeyFromHex(secret)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create symmetric key from provided secret: %w",
			err,
		)
	}

	p := paseto.NewParserWithoutExpiryCheck()
	tkn, err := p.ParseV4Local(key, enc, implicit)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt token: %w", err,
		)
	}

	return &token{Token: *tkn}, nil
}
