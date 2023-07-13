package user

type GetProfileQuery struct {
	UserID string `form:"user_id" binding:"required"`
}

type SignInPayload struct {
	Type int    `form:"type" binding:"required"`
	Code string `form:"code" binding:"required"`
}

type GitHubTokenPayload struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type GitHubTokenSuccessResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

type GitHubTokenResp struct {
	GitHubTokenSuccessResp
	GithubTokenError
}
