package recaptcha

import (
	"context"
	"errors"
	"net"
)

var ErrInvalidCaptcha = errors.New("invalid recaptcha")

type Validator interface {
	Validate(
		ctx context.Context,
		token string,
		clientIP net.IP,
	) (bool, error)
}
