package wallpaper

type SearchQuery struct {
	Keywords    string `form:"q"`
	Categories  string `form:"categories" binding:"required"`
	Purity      string `form:"purity" binding:"required"`
	Sorting     string `form:"sorting" binding:"required"`
	Order       string `form:"order" binding:"required"`
	TopRange    string `form:"topRange"`
	AtLeast     string `form:"atleast"`
	Resolutions string `form:"resolutions"`
	Ratios      string `form:"ratios"`
	Colors      string `form:"colors"`
	AIArtFilter string `form:"ai_art_filter" binding:"required"`
	Page        string `form:"page" binding:"required"`
}

type GetInfoQuery struct {
	URL string `form:"url" binding:"required"`
}

type SearchResp struct {
	Data []Wallpaper `json:"data"`
}

type WallpaperInfoResp struct {
	Data Wallpaper `json:"data"`
}

type SearchAPIRes struct {
	Error  error
	Result SearchResp
}
