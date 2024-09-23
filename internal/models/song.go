package models

type Song struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumCover  string `json:"album_cover"`
	Duration    int    `json:"duration"`
	ReleaseDate string `json:"release_date"`
	Popularity  int    `json:"popularity"`
	PreviewURL  string `json:"preview_url"`
	ExternalURL string `json:"external_url"`
}
