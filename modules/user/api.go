package user

type GetProfileQuery struct {
	UserID string `form:"user_id"`
}

type UpdateProfilePayload struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}

type LoginPayload struct {
	Type int    `json:"type" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type LoginRes struct {
	User
	Token string `json:"token"`
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

type GithubProfileResp struct {
	Login                   string `json:"login"`
	ID                      int    `json:"id"`
	NodeID                  string `json:"node_id"`
	AvatarURL               string `json:"avatar_url"`
	GravatarID              string `json:"gravatar_id"`
	URL                     string `json:"url"`
	HTMLURL                 string `json:"html_url"`
	FollowersURL            string `json:"followers_url"`
	FollowingURL            string `json:"following_url"`
	GistsURL                string `json:"gists_url"`
	StarredURL              string `json:"starred_url"`
	SubscriptionsURL        string `json:"subscriptions_url"`
	OrganizationsURL        string `json:"organizations_url"`
	ReposURL                string `json:"repos_url"`
	EventsURL               string `json:"events_url"`
	ReceivedEventsURL       string `json:"received_events_url"`
	Type                    string `json:"type"`
	SiteAdmin               bool   `json:"site_admin"`
	Name                    string `json:"name"`
	Company                 string `json:"company"`
	Blog                    string `json:"blog"`
	Location                string `json:"location"`
	Email                   string `json:"email"`
	Hireable                bool   `json:"hireable"`
	Bio                     string `json:"bio"`
	TwitterUsername         string `json:"twitter_username"`
	PublicRepos             int    `json:"public_repos"`
	PublicGists             int    `json:"public_gists"`
	Followers               int    `json:"followers"`
	Following               int    `json:"following"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
	PrivateGists            int    `json:"private_gists"`
	TotalPrivateRepos       int    `json:"total_private_repos"`
	OwnedPrivateRepos       int    `json:"owned_private_repos"`
	DiskUsage               int    `json:"disk_usage"`
	Collaborators           int    `json:"collaborators"`
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`
}

type YDWechatQRCodeResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		QRURL      string `json:"qrUrl"`
		TempUserID string `json:"tempUserId"`
	} `json:"data"`
}

type GetYDQRCodeResp struct {
	QRCodeURL  string `json:"qrcode_url"`
	TempUserID string `json:"temp_user_id"`
}

type GetYDLoginStatusQuery struct {
	TempUserID string `form:"temp_user_id" binding:"required"`
}

type GetYDLoginStatusResp struct {
	HasLogin bool `json:"has_login"`
}

type WeChatUserInfo struct {
	OpenID    string `json:"openId"`
	Nickname  string `json:"nickName"`
	Gender    string `json:"gender"`
	AvatarUrl string `json:"avatarUrl"`
}

type YDCallbackPayload struct {
	TempUserId   string         `json:"tempUserId" binding:"required"`
	ScanSuccess  bool           `json:"scanSuccess"`
	CancelLogin  bool           `json:"cancelLogin"`
	WxMaUserInfo WeChatUserInfo `json:"wxMaUserInfo"`
}
