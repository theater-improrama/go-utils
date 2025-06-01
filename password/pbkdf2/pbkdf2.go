package pbkdf2

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"hash"

	jsonext "github.com/theater-improrama/go-utils/jsonext"
	"github.com/theater-improrama/go-utils/password"
	"golang.org/x/crypto/pbkdf2"
)

func Provide(sG password.SaltGenerator, p Params) password.ProviderFn {
	return func() (password.Provider, error) {
		s, err := sG.Generate()
		if err != nil {
			return nil, err
		}

		return NewSha512(
			s,
			p,
		), nil
	}
}

const TypeSha512 = "pbkdf2-sha512"

type Pbkdf2Sha512HasherValidator struct {
	hFn    func() hash.Hash
	salt   []byte
	params Params
}

type Params struct {
	KeyLen int `json:"key_len"`
	Iter   int `json:"iter"`
}

type encodableHash struct {
	Derivative jsonext.Base64Arr `json:"derivative"`
	Salt       jsonext.Base64Arr `json:"salt"`
	Params     Params            `json:"params"`
}

func Decode(eH string) ([]byte, []byte, Params, error) {
	bs, err := base64.StdEncoding.DecodeString(eH)
	if err != nil {
		return nil, nil, Params{}, err
	}

	var h encodableHash
	if err := json.Unmarshal(bs, &h); err != nil {
		return nil, nil, Params{}, err
	}

	return h.Derivative, h.Salt, h.Params, nil
}

func Encode(dK []byte, salt []byte, p Params) (string, error) {
	eH := encodableHash{
		Derivative: dK,
		Salt:       salt,
		Params:     p,
	}

	s, err := json.Marshal(&eH)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(s), nil
}

func NewSha512(
	salt []byte,
	p Params,
) *Pbkdf2Sha512HasherValidator {
	return &Pbkdf2Sha512HasherValidator{
		hFn:    sha512.New,
		salt:   salt,
		params: p,
	}
}

func (p *Pbkdf2Sha512HasherValidator) Type() string {
	return TypeSha512
}

func (p *Pbkdf2Sha512HasherValidator) HashEncode(password []byte) (string, error) {
	dK, err := p.Hash(password)
	if err != nil {
		return "", err
	}

	return Encode(dK, p.salt, p.params)
}

func (p *Pbkdf2Sha512HasherValidator) hash(password []byte, salt []byte) ([]byte, error) {
	dK := pbkdf2.Key(
		password,
		salt,
		p.params.Iter,
		p.params.KeyLen,
		p.hFn,
	)

	return dK, nil
}

func (p *Pbkdf2Sha512HasherValidator) Hash(password []byte) ([]byte, error) {
	return p.hash(password, p.salt)
}

func (p *Pbkdf2Sha512HasherValidator) Validate(dk []byte, password []byte) (bool, error) {
	pDk, err := p.Hash(password)
	if err != nil {
		return false, err
	}

	if bytes.Compare(dk, pDk) != 0 {
		return false, nil
	}

	return true, nil
}

var _ password.Provider = (*Pbkdf2Sha512HasherValidator)(nil)
