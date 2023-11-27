package model

type Wallpaper struct {
	WallpaperID string   `json:"id" bson:"id"`
	URL         string   `json:"url" bson:"url"`
	ShortURL    string   `json:"short_url" bson:"short_url"`
	Views       int      `json:"views" bson:"views"`
	Purity      string   `json:"purity" bson:"purity"`
	Category    string   `json:"category" bson:"category"`
	DimensionX  int      `json:"dimension_x" bson:"dimension_x"`
	DimensionY  int      `json:"dimension_y" bson:"dimension_y"`
	Ratio       string   `json:"ratio" bson:"ratio"`
	Resolution  string   `json:"resolution" bson:"resolution"`
	File_size   int      `json:"file_size" bson:"file_size"`
	File_type   string   `json:"file_type" bson:"file_type"`
	Created_at  string   `json:"created_at" bson:"created_at"`
	Colors      []string `json:"colors" bson:"colors"`
	Path        string   `json:"path" bson:"path"`
	Thumbs      Thumbs   `json:"thumbs" bson:"thumbs"`
	Tags        []Tag    `json:"tags" bson:"tags"`
}

type Tag struct {
	TagID       int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Alias       string `json:"alias" bson:"alias"`
	Category_id int    `json:"category_id" bson:"category_id"`
	Category    string `json:"category" bson:"category"`
	Purity      string `json:"purity" bson:"purity"`
	Created_at  string `json:"created_at" bson:"created_at"`
}

type Thumbs struct {
	Large    string `json:"large" bson:"large"`
	Original string `json:"original" bson:"original"`
	Small    string `json:"small" bson:"small"`
}
