package internal

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Ensure Password implements the proper pgtypes interfaces so it can be
	// used with pgx seamlessly.
	_ pgtype.BytesScanner = (*Password)(nil)
	_ pgtype.BytesValuer  = (*Password)(nil)
)

type Password struct {
	plaintext *string
	hash      []byte
}

func NewPasswordFromPlaintext(password string) (*Password, error) {
	p := &Password{}

	err := p.Set(password)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plaintext = &password
	p.hash = hash

	return nil
}

func (p *Password) Matches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (p *Password) ScanBytes(v []byte) error {
	newVal := make([]byte, len(v))
	copy(newVal, v)
	p.hash = newVal
	return nil
}

func (p *Password) BytesValue() ([]byte, error) {
	return p.hash, nil
}
