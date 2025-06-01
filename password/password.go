package password

type Typer interface {
	Type() string
}

type HashEncoder interface {
	Typer
	HashEncode(password []byte) (string, error)
}

type Hasher interface {
	Typer
	Hash(password []byte) ([]byte, error)
}

type Validator interface {
	Typer
	Validate(dK []byte, password []byte) (bool, error)
}

type Provider interface {
	Typer
	HashEncoder
	Hasher
	Validator
}

type SaltGenerator interface {
	Generate() ([]byte, error)
}
