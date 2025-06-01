package argon2id

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"

	"github.com/theater-improrama/go-utils/jsonext"
	"github.com/theater-improrama/go-utils/password"
	"golang.org/x/crypto/argon2"
)

const Type = "argon2id"

func Provide(sG password.SaltGenerator, p Params) password.ProviderFn {
	return func() (password.Provider, error) {
		s, err := sG.Generate()
		if err != nil {
			return nil, err
		}

		return New(
			s,
			p,
		), nil
	}
}

type Argon2idHasherValidator struct {
	salt   []byte
	params Params
}

type Params struct {
	Time    uint32 `json:"time"`
	Memory  uint32 `json:"memory"`
	Threads uint8  `json:"threads"`
	KeyLen  uint32 `json:"key_len"`
}

func New(
	salt []byte,
	params Params,
) *Argon2idHasherValidator {
	return &Argon2idHasherValidator{
		params: params,
		salt:   salt,
	}
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

func (v *Argon2idHasherValidator) HashEncode(password []byte) (string, error) {
	dK, err := v.hash(password, v.salt)
	if err != nil {
		return "", err
	}

	return Encode(dK, v.salt, v.params)
}

func (v *Argon2idHasherValidator) Type() string {
	return Type
}

func (v *Argon2idHasherValidator) hash(password []byte, salt []byte) ([]byte, error) {
	dK := argon2.IDKey(
		password,
		salt,
		v.params.Time,
		v.params.Memory,
		v.params.Threads,
		v.params.KeyLen,
	)

	return dK, nil
}

func (v *Argon2idHasherValidator) Hash(password []byte) ([]byte, error) {
	return v.hash(password, v.salt)
}

func (v *Argon2idHasherValidator) Validate(dK []byte, password []byte) (bool, error) {
	pDk, err := v.Hash(password)
	if err != nil {
		return false, err
	}

	kL := int32(len(dK))
	pKL := int32(len(pDk))

	if subtle.ConstantTimeEq(kL, pKL) == 0 {
		return false, nil
	}

	if subtle.ConstantTimeCompare(dK, pDk) != 1 {
		return false, nil
	}

	return true, nil
}

var _ password.Provider = (*Argon2idHasherValidator)(nil)
