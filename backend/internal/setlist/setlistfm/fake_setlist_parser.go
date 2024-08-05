package setlistfm

import (
	"festwrap/internal/setlist"
)

type FakeSetlistParser struct {
	parseArgs []byte
	response  *setlist.Setlist
	err       error
}

func (p *FakeSetlistParser) GetParseArgs() []byte {
	return p.parseArgs
}

func (p *FakeSetlistParser) SetError(err error) {
	p.err = err
}

func (p *FakeSetlistParser) SetReponse(response *setlist.Setlist) {
	p.response = response
}

func (p *FakeSetlistParser) Parse(setlist []byte) (*setlist.Setlist, error) {
	p.parseArgs = setlist
	if p.err != nil {
		return nil, p.err
	}

	return p.response, nil
}
