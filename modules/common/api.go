package common

type GetBucketTokenResp struct {
	Token string `json:"token"`
}

type GetSignedURLQuery struct {
	Key string `form:"key" binding:"required"`
}

type GetSignedURLResp struct {
	URL string `json:"url"`
}
