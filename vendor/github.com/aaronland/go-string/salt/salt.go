package salt

import (
	"errors"
	"github.com/aaronland/go-string/random"
)

var min_length int

func init() {
	min_length = 32
}

type Salt struct {
	salt string
}

func (s *Salt) String() string {
	return s.salt
}

func (s *Salt) Bytes() []byte {
	return []byte(s.salt)
}

type SaltOptions struct {
	Length int
	ASCII  bool
}

func DefaultSaltOptions() *SaltOptions {

	opts := SaltOptions{
		Length: min_length,
		ASCII:  false,
	}

	return &opts
}

func NewRandomSalt(opts *SaltOptions) (*Salt, error) {

	str_opts := random.DefaultOptions()
	str_opts.Length = opts.Length
	str_opts.ASCII = opts.ASCII

	s, err := random.String(str_opts)

	if err != nil {
		return nil, err
	}

	return NewSaltFromString(s)
}

func NewSaltFromString(s string) (*Salt, error) {

	_, err := IsValidSalt(s)

	salt := Salt{s}

	if err != nil {
		return nil, err
	}

	return &salt, nil
}

func IsValidSalt(s string) (bool, error) {

	if len(s) < min_length {
		return false, errors.New("salt is too short")
	}

	return true, nil
}
