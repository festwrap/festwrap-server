package tidal

type tidalTrack struct {
	Id string `json:"id"`
}

type tidalResponse struct {
	Results []tidalTrack `json:"included"`
}
