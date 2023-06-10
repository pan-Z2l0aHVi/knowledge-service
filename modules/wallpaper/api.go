package wallpaper

type SearchParams struct {
	Categories  string `form:"categories" binding:"required"`
	Purity      string `form:"purity" binding:"required"`
	Sorting     string `form:"sorting" binding:"required"`
	Order       string `form:"order" binding:"required"`
	TopRange    string `form:"topRange"`
	AtLeast     string `form:"atleast"`
	Resolutions string `form:"resolutions"`
	Ratios      string `form:"ratios"`
	Colors      string `form:"colors"`
	Page        string `form:"page" binding:"required"`
	AIArtFilter string `form:"ai_art_filter" binding:"required"`
}

type GetInfoParams struct {
	URL string `form:"url" binding:"required"`
}

type WallhavenSearchResp struct {
	Data []Wallpaper `json:"data"`
}

type WallhavenInfoResp struct {
	Data Wallpaper `json:"data"`
}
