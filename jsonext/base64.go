package jsonext

import (
	"encoding/base64"
	"encoding/json"
)

type Base64Arr []byte

func (b *Base64Arr) MarshalJSON() ([]byte, error) {
	bs := make([]byte, base64.StdEncoding.EncodedLen(len(*b)))

	base64.StdEncoding.Encode(bs, *b)

	return json.Marshal(bs)
}

func (b *Base64Arr) UnmarshalJSON(data []byte) error {
	var jBs []byte
	if err := json.Unmarshal(data, &jBs); err != nil {
		return err
	}

	bs := make([]byte, base64.StdEncoding.DecodedLen(len(jBs)))
	n, err := base64.StdEncoding.Decode(bs, jBs)
	if err != nil {
		return err
	}

	*b = bs[:n]
	return nil
}
