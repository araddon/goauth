package oauth

import ()

type RequestToken struct {
	Token    string `json:"token"`
	Secret   string `json:"secret"`
	Verifier string `json:"verifier"`
}

type AccessToken struct {
	Id       string `json:"id"`
	Token    string `json:"token"`
	Secret   string `json:"secret"`
	UserRef  string `json:"userref"`
	Verifier string `json:"verifier"`
	Service  string `json:"service"`
}
