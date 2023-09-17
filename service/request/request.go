package request

type PostGetOAuthTokenRequest struct {
}

type PostGetAccessTokenRequest struct {
	OAuthToken    string `form:"oauth_token" binding:"required"`
	OAuthVerifier string `form:"oauth_verifier" binding:"required"`
	Address       string `form:"address" binding:"required"`
	SignedMessage string `form:"signed_message" binding:"required"`
}
