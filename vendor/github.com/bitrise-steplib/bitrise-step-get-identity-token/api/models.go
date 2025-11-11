package api

type GetIdentityTokenParameter struct {
	Subject  string `json:"sub"`
	Audience string `json:"aud"`
}

type GetIdentityTokenResponse struct {
	Token     string `json:"id_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

type errorReponse struct {
	Message string `json:"error_msg"`
}
