package token_test

import (
	"testing"

	"github.com/murtaza-u/ddos/token"

	"aidanwoods.dev/go-paseto"
)

func TestDecrypt(t *testing.T) {
	tkn, err := token.New(token.Params{
		Body: map[string]any{"foo": "bar"},
	})
	if err != nil {
		t.Errorf("NewToken: %v", err)
	}

	key := paseto.NewV4SymmetricKey()
	enc, err := tkn.Encrypt(key.ExportHex(), nil)
	if err != nil {
		t.Errorf("*token.Encrypt: %v", err)
	}

	tkn, err = token.Decrypt(enc, key.ExportHex(), nil)
	if err != nil {
		t.Errorf("*token.Decrypt: %v", err)
	}

	bar, err := tkn.GetString("foo")
	if bar != "bar" {
		t.Error("incorrectly decrypted token")
	}
}
