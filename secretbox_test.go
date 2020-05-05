package secretbox

import (
	"testing"
)

func TestSecretbox(t *testing.T) {

	secret := "s33kret"
	salt := "s4lty"
	plain := "hello world"

	opts := NewSecretboxOptions()
	opts.Salt = salt

	sb, err := NewSecretbox(secret, opts)

	if err != nil {
		t.Fatal(err)
	}

	locked, err := sb.Lock([]byte(plain))

	if err != nil {
		t.Fatal(err)
	}

	unlocked, err := sb.Unlock(locked)

	if err != nil {
		t.Fatal(err)
	}

	if string(unlocked.String()) != plain {
		t.Fatal("Unlock failed")
	}
}
