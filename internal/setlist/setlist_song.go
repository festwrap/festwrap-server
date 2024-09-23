package setlist

type Song struct {
	title string
}

func NewSong(title string) Song {
	return Song{title: title}
}

func (s Song) GetTitle() string {
	return s.title
}
