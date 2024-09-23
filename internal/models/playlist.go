package models

type Playlist struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
	Songs  []Song `json:"songs,omitempty"`
}

type CreatePlaylistDto struct {
	Name string `json:"name"`
}

type UpdatePlaylistDto struct {
	Name string `json:"name"`
}
