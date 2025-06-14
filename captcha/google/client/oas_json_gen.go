// Code generated by ogen, DO NOT EDIT.

package client

import (
	"math/bits"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/ogen-go/ogen/validate"
)

// Encode encodes SiteverifyResponseData as json.
func (o OptSiteverifyResponseData) Encode(e *jx.Encoder) {
	if !o.Set {
		return
	}
	o.Value.Encode(e)
}

// Decode decodes SiteverifyResponseData from json.
func (o *OptSiteverifyResponseData) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New("invalid: unable to decode OptSiteverifyResponseData to nil")
	}
	o.Set = true
	if err := o.Value.Decode(d); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s OptSiteverifyResponseData) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *OptSiteverifyResponseData) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *SiteverifyResponse) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *SiteverifyResponse) encodeFields(e *jx.Encoder) {
	{
		if s.Data.Set {
			e.FieldStart("data")
			s.Data.Encode(e)
		}
	}
}

var jsonFieldsNameOfSiteverifyResponse = [1]string{
	0: "data",
}

// Decode decodes SiteverifyResponse from json.
func (s *SiteverifyResponse) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode SiteverifyResponse to nil")
	}

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "data":
			if err := func() error {
				s.Data.Reset()
				if err := s.Data.Decode(d); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"data\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode SiteverifyResponse")
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *SiteverifyResponse) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *SiteverifyResponse) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}

// Encode implements json.Marshaler.
func (s *SiteverifyResponseData) Encode(e *jx.Encoder) {
	e.ObjStart()
	s.encodeFields(e)
	e.ObjEnd()
}

// encodeFields encodes fields.
func (s *SiteverifyResponseData) encodeFields(e *jx.Encoder) {
	{
		e.FieldStart("success")
		e.Bool(s.Success)
	}
	{
		e.FieldStart("challenge_ts")
		e.Str(s.ChallengeTs)
	}
	{
		e.FieldStart("hostname")
		e.Str(s.Hostname)
	}
	{
		e.FieldStart("error-codes")
		e.ArrStart()
		for _, elem := range s.ErrorMinusCodes {
			e.Str(elem)
		}
		e.ArrEnd()
	}
}

var jsonFieldsNameOfSiteverifyResponseData = [4]string{
	0: "success",
	1: "challenge_ts",
	2: "hostname",
	3: "error-codes",
}

// Decode decodes SiteverifyResponseData from json.
func (s *SiteverifyResponseData) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New("invalid: unable to decode SiteverifyResponseData to nil")
	}
	var requiredBitSet [1]uint8

	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "success":
			requiredBitSet[0] |= 1 << 0
			if err := func() error {
				v, err := d.Bool()
				s.Success = bool(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"success\"")
			}
		case "challenge_ts":
			requiredBitSet[0] |= 1 << 1
			if err := func() error {
				v, err := d.Str()
				s.ChallengeTs = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"challenge_ts\"")
			}
		case "hostname":
			requiredBitSet[0] |= 1 << 2
			if err := func() error {
				v, err := d.Str()
				s.Hostname = string(v)
				if err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"hostname\"")
			}
		case "error-codes":
			requiredBitSet[0] |= 1 << 3
			if err := func() error {
				s.ErrorMinusCodes = make([]string, 0)
				if err := d.Arr(func(d *jx.Decoder) error {
					var elem string
					v, err := d.Str()
					elem = string(v)
					if err != nil {
						return err
					}
					s.ErrorMinusCodes = append(s.ErrorMinusCodes, elem)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return errors.Wrap(err, "decode field \"error-codes\"")
			}
		default:
			return d.Skip()
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "decode SiteverifyResponseData")
	}
	// Validate required fields.
	var failures []validate.FieldError
	for i, mask := range [1]uint8{
		0b00001111,
	} {
		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
			// Mask only required fields and check equality to mask using XOR.
			//
			// If XOR result is not zero, result is not equal to expected, so some fields are missed.
			// Bits of fields which would be set are actually bits of missed fields.
			missed := bits.OnesCount8(result)
			for bitN := 0; bitN < missed; bitN++ {
				bitIdx := bits.TrailingZeros8(result)
				fieldIdx := i*8 + bitIdx
				var name string
				if fieldIdx < len(jsonFieldsNameOfSiteverifyResponseData) {
					name = jsonFieldsNameOfSiteverifyResponseData[fieldIdx]
				} else {
					name = strconv.Itoa(fieldIdx)
				}
				failures = append(failures, validate.FieldError{
					Name:  name,
					Error: validate.ErrFieldRequired,
				})
				// Reset bit.
				result &^= 1 << bitIdx
			}
		}
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}

	return nil
}

// MarshalJSON implements stdjson.Marshaler.
func (s *SiteverifyResponseData) MarshalJSON() ([]byte, error) {
	e := jx.Encoder{}
	s.Encode(&e)
	return e.Bytes(), nil
}

// UnmarshalJSON implements stdjson.Unmarshaler.
func (s *SiteverifyResponseData) UnmarshalJSON(data []byte) error {
	d := jx.DecodeBytes(data)
	return s.Decode(d)
}
