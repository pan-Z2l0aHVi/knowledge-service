package entity

type GetBucketTokenResp struct {
	Token string `json:"token"`
}

type GetSignedURLQuery struct {
	Key string `form:"key" binding:"required"`
}

type GetSignedURLResp struct {
	URL string `json:"url"`
}

type ReportPayload struct {
	Data  []interface{} `json:"data" binding:"required"`
	Token string        `json:"token"`
}

type GetStaticsQuery struct {
	StartTimestamp int64 `form:"start_timestamp"`
	EndTimestamp   int64 `form:"end_timestamp"`
}

type GetStaticsResp struct {
	PV int64 `json:"pv"`
	UV int64 `json:"uv"`
}
