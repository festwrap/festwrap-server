package setlist

type Song struct {
	title string
}

func (s Song) GetTitle() string {
	return s.title
}

func NewSong(title string) Song {
	return Song{title: title}
}
