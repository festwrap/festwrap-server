package artist

type Artist struct {
	Name     string `json:"name"`
	ImageUri string `json:"imageUri,omitempty"`
}

func NewArtist(name string) Artist {
	return Artist{Name: name}
}

func NewArtistWithImageUri(name string, imageUri string) Artist {
	return Artist{Name: name, ImageUri: imageUri}
}

func (a *Artist) SetImageUri(imageUri string) {
	a.ImageUri = imageUri
}
