// Code generated by ogen, DO NOT EDIT.

package client

// NewOptSiteverifyResponseData returns new OptSiteverifyResponseData with value set to v.
func NewOptSiteverifyResponseData(v SiteverifyResponseData) OptSiteverifyResponseData {
	return OptSiteverifyResponseData{
		Value: v,
		Set:   true,
	}
}

// OptSiteverifyResponseData is optional SiteverifyResponseData.
type OptSiteverifyResponseData struct {
	Value SiteverifyResponseData
	Set   bool
}

// IsSet returns true if OptSiteverifyResponseData was set.
func (o OptSiteverifyResponseData) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptSiteverifyResponseData) Reset() {
	var v SiteverifyResponseData
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptSiteverifyResponseData) SetTo(v SiteverifyResponseData) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptSiteverifyResponseData) Get() (v SiteverifyResponseData, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptSiteverifyResponseData) Or(d SiteverifyResponseData) SiteverifyResponseData {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/SiteverifyForm
type SiteverifyForm struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	Remoteip string `json:"remoteip"`
}

// GetSecret returns the value of Secret.
func (s *SiteverifyForm) GetSecret() string {
	return s.Secret
}

// GetResponse returns the value of Response.
func (s *SiteverifyForm) GetResponse() string {
	return s.Response
}

// GetRemoteip returns the value of Remoteip.
func (s *SiteverifyForm) GetRemoteip() string {
	return s.Remoteip
}

// SetSecret sets the value of Secret.
func (s *SiteverifyForm) SetSecret(val string) {
	s.Secret = val
}

// SetResponse sets the value of Response.
func (s *SiteverifyForm) SetResponse(val string) {
	s.Response = val
}

// SetRemoteip sets the value of Remoteip.
func (s *SiteverifyForm) SetRemoteip(val string) {
	s.Remoteip = val
}

// Ref: #/components/schemas/SiteverifyResponse
type SiteverifyResponse struct {
	Data OptSiteverifyResponseData `json:"data"`
}

// GetData returns the value of Data.
func (s *SiteverifyResponse) GetData() OptSiteverifyResponseData {
	return s.Data
}

// SetData sets the value of Data.
func (s *SiteverifyResponse) SetData(val OptSiteverifyResponseData) {
	s.Data = val
}

type SiteverifyResponseData struct {
	Success         bool     `json:"success"`
	ChallengeTs     string   `json:"challenge_ts"`
	Hostname        string   `json:"hostname"`
	ErrorMinusCodes []string `json:"error-codes"`
}

// GetSuccess returns the value of Success.
func (s *SiteverifyResponseData) GetSuccess() bool {
	return s.Success
}

// GetChallengeTs returns the value of ChallengeTs.
func (s *SiteverifyResponseData) GetChallengeTs() string {
	return s.ChallengeTs
}

// GetHostname returns the value of Hostname.
func (s *SiteverifyResponseData) GetHostname() string {
	return s.Hostname
}

// GetErrorMinusCodes returns the value of ErrorMinusCodes.
func (s *SiteverifyResponseData) GetErrorMinusCodes() []string {
	return s.ErrorMinusCodes
}

// SetSuccess sets the value of Success.
func (s *SiteverifyResponseData) SetSuccess(val bool) {
	s.Success = val
}

// SetChallengeTs sets the value of ChallengeTs.
func (s *SiteverifyResponseData) SetChallengeTs(val string) {
	s.ChallengeTs = val
}

// SetHostname sets the value of Hostname.
func (s *SiteverifyResponseData) SetHostname(val string) {
	s.Hostname = val
}

// SetErrorMinusCodes sets the value of ErrorMinusCodes.
func (s *SiteverifyResponseData) SetErrorMinusCodes(val []string) {
	s.ErrorMinusCodes = val
}
