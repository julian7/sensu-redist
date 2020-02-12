package redist

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julian7/sensulib/sensuasset"
	"gopkg.in/yaml.v2"
)

// Download copies asset items locally, preparing for relocation
func (r *Redist) Download(source *url.URL) error {
	asset, err := r.downloadItem(source, "asset file")
	if err != nil {
		return err
	}

	r.Asset = asset

	if err := r.parseAsset(); err != nil {
		return err
	}

	if err := r.downloadBuilds(); err != nil {
		return err
	}

	return nil
}

func (r *Redist) downloadItem(source *url.URL, descr string) (*Asset, error) {
	basename := path.Base(source.Path)
	local := filepath.Join(r.Dir, basename)

	if descr == "" {
		descr = basename
	}

	if err := downloadFile(local, source.String(), descr); err != nil {
		return nil, fmt.Errorf("downloading %s: %w", descr, err)
	}

	return &Asset{Desc: descr, Local: local, Orig: source}, nil
}

func (r *Redist) parseAsset() error {
	r.AssetSpec = &sensuasset.AssetSpec{}

	contents, err := ioutil.ReadFile(r.Asset.Local)
	if err != nil {
		return fmt.Errorf("loading asset file: %w", err)
	}

	if err := yaml.Unmarshal(contents, &r.AssetSpec); err != nil {
		return fmt.Errorf("loading asset file: %w", err)
	}

	return nil
}

func (r *Redist) downloadBuilds() error {
	r.Builds = []*Asset{}

	for _, build := range r.AssetSpec.Spec.Builds {
		url, err := url.Parse(build.URL)
		if err != nil {
			return fmt.Errorf("parsing build URL: %w", err)
		}

		asset, err := r.downloadItem(url, "")
		if err != nil {
			return err
		}

		asset.BuildSpec = build

		r.Builds = append(r.Builds, asset)
	}

	return nil
}

func downloadFile(local, url, descr string) error {
	st, err := os.Stat(local)
	if err == nil {
		if st.Mode().IsRegular() {
			return nil
		}

		return fmt.Errorf("%q already exists", local)
	}

	cl := &http.Client{Timeout: 10 * time.Second}

	resp, err := cl.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 OK while downloading %s, got %s", descr, resp.Status)
	}

	f, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("opening temp file to write %s: %w", descr, err)
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("writing asset file: %w", err)
	}

	return nil
}
