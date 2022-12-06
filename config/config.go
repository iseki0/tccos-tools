package config

import (
	"encoding/base64"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
	"tccos-tools/optstr"
)

type C struct {
	BaseURL *cos.BaseURL
	Auth    *cos.AuthorizationTransport
}

//go:generate stringer  -type configError -output config_error_string.go -linecomment
type configError int

const (
	_ configError = iota
	errConfigUnused
)

func (i configError) Error() string {
	return i.String()
}

func Parse(input string) (*C, error) {
	data, e := base64.StdEncoding.DecodeString(input)
	if e != nil {
		return nil, fmt.Errorf("option string decode base64: %w", e)
	}
	options, e := optstr.ParseString(string(data))
	if e != nil {
		return nil, fmt.Errorf("parse options: %w", e)
	}
	var m = make(map[string]string)
	for _, it := range options {
		m[it.Key] = it.Value
	}

	var c C
	c.BaseURL = &cos.BaseURL{}
	c.Auth = &cos.AuthorizationTransport{}

	c.BaseURL.ServiceURL, e = url.Parse(m["service"])
	if e != nil {
		return nil, fmt.Errorf("service url: %w", e)
	}
	c.BaseURL.BucketURL, e = url.Parse(m["bucket"])
	if e != nil {
		return nil, fmt.Errorf("bucket url: %w", e)
	}
	c.Auth.SecretID = m["sid"]
	c.Auth.SecretKey = m["sk"]
	return &c, nil
}
