package redist

import (
	"net/url"

	"github.com/julian7/sensulib/sensuasset"
)

type Redist struct {
	*sensuasset.AssetSpec
	Asset  *Asset
	Builds []*Asset
	Dir    string
}

type Asset struct {
	*sensuasset.BuildSpec
	Desc  string
	Local string
	New   *url.URL
	Orig  *url.URL
}

func New(dir string) (*Redist, error) {
	return &Redist{Dir: dir}, nil
}
