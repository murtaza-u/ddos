package token_test

import (
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/murtaza-u/ddos/token"
)

func TestNew(t *testing.T) {
	body := make(map[string]any)
	body["foo"] = "bar"

	var err error

	// invalid `nbf` field
	_, err = token.New(token.Params{
		Body:      body,
		Expiry:    time.Now().Add(time.Minute * 5),
		NotBefore: time.Now().Add(time.Minute * 10),
	})
	if err != token.ErrInvalidNB {
		t.Errorf(
			"NewToken: expected: %v | got: %v",
			token.ErrInvalidNB, err,
		)
	}

	// body = nil
	_, err = token.New(token.Params{
		Body: nil,
	})
	if err != token.ErrMissingBody {
		t.Errorf(
			"NewToken: expected: %v | got: %v",
			token.ErrMissingBody, err,
		)
	}

	// body = empty map
	_, err = token.New(token.Params{
		Body: make(map[string]any),
	})
	if err != token.ErrMissingBody {
		t.Errorf(
			"NewToken: expected: %v | got: %v",
			token.ErrMissingBody, err,
		)
	}

	// no errors
	_, err = token.New(token.Params{
		Body: body,
	})
	if err != nil {
		t.Errorf(
			"NewToken: expected: nil | got: %v", err,
		)
	}
}

func TestEncrypt(t *testing.T) {
	p := token.Params{
		Body: map[string]any{"foo": "bar"},
	}

	tkn, err := token.New(p)
	if err != nil {
		t.Errorf("NewToken: %v", err)
	}

	// invalid secret
	_, err = tkn.Encrypt("thiswillfail", nil)
	if err == nil {
		t.Error("*token.Encrypt: expected: an error | got: nil")
	}

	// valid secret. Should pass.
	key := paseto.NewV4SymmetricKey()
	enc, err := tkn.Encrypt(key.ExportHex(), nil)
	if err != nil {
		t.Errorf("*token.Encrypt: %v", err)
	}
	t.Logf("encrypted token: %s\n", enc)
}

func TestIsExpired(t *testing.T) {
	p := token.Params{
		Expiry: time.Now().Add(time.Millisecond),
		Body:   map[string]any{"foo": "bar"},
	}

	tkn, err := token.New(p)
	if err != nil {
		t.Errorf("NewToken: %v", err)
	}

	time.Sleep(time.Millisecond)

	if !tkn.IsExpired() {
		t.Errorf("*token.IsExpired: invalid outcome")
	}
}

func TestRefresh(t *testing.T) {
	p := token.Params{
		Expiry: time.Now().Add(time.Millisecond),
		Body:   map[string]any{"foo": "bar"},
	}

	tkn, err := token.New(p)
	if err != nil {
		t.Errorf("NewToken: %v", err)
	}

	time.Sleep(time.Millisecond)

	err = tkn.Refresh()
	if err != nil {
		t.Errorf("Refresh: %v", err)
	}
}
