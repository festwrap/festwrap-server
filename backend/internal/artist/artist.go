package artist

type Artist struct {
	name     string
	imageUri string
}

func (a *Artist) SetImageUri(imageUri string) {
	a.imageUri = imageUri
}

func NewArtist(name string) Artist {
	return Artist{name: name}
}

func NewArtistWithImageUri(name string, imageUri string) Artist {
	return Artist{name: name, imageUri: imageUri}
}
