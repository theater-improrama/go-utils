package google

import (
	"context"
	"net"

	recaptcha "github.com/theater-improrama/go-utils/captcha"
	"github.com/theater-improrama/go-utils/captcha/google/client"
)

type recaptchaValidator struct {
	c      *client.Client
	secret string
}

func NewValidator(secret string) (recaptcha.Validator, error) {
	c, err := client.NewClient("https://www.google.com/recaptcha")
	if err != nil {
		return nil, err
	}

	return &recaptchaValidator{
		c:      c,
		secret: secret,
	}, nil
}

func (r *recaptchaValidator) Validate(
	ctx context.Context,
	token string,
	clientIP net.IP,
) (bool, error) {
	resp, err := r.c.Siteverify(ctx, &client.SiteverifyForm{
		Response: token,
		Secret:   r.secret,
		Remoteip: clientIP.String(),
	})
	if err != nil {
		return false, err
	}

	if err := resp.Validate(); err != nil {
		return false, err
	}

	if resp.Data.Set && len(resp.Data.Value.ErrorMinusCodes) > 0 {
		return false, nil
	}

	return true, nil
}

var _ recaptcha.Validator = (*recaptchaValidator)(nil)
