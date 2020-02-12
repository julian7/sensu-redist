package redist

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"

	"gopkg.in/yaml.v3"
)

func (r *Redist) RelocateTo(target *url.URL) error {
	if err := r.relocateURL(target, r.Asset); err != nil {
		return err
	}

	for _, build := range r.Builds {
		if err := r.relocateURL(target, build); err != nil {
			return err
		}
	}

	out, err := yaml.Marshal(r.AssetSpec)
	if err != nil {
		return fmt.Errorf("converting modified asset to YAML: %w", err)
	}

	if err := ioutil.WriteFile(r.Asset.Local, out, 0644); err != nil {
		return fmt.Errorf("writing modified asset file: %w", err)
	}

	return nil
}

func (r *Redist) relocateURL(target *url.URL, asset *Asset) error {
	base, err := url.Parse(path.Base(asset.Orig.Path))
	if err != nil {
		return fmt.Errorf("relocating %s: %w", asset.Desc, err)
	}

	newURL := target.ResolveReference(base)
	asset.New = newURL

	if asset.BuildSpec != nil {
		asset.BuildSpec.URL = newURL.String()
	}

	return nil
}
