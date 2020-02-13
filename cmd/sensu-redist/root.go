package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/julian7/sensu-redist/redist"
	"github.com/spf13/cobra"
)

type redistConfig struct {
	Source  string
	srcURL  *url.URL
	Target  string
	tarURL  *url.URL
	Upload  string
	Workdir string
}

func rootCmd() (*cobra.Command, error) {
	conf := &redistConfig{}
	cmd := &cobra.Command{
		Use:     "sensu-redist NAME",
		Short:   "Downloads and re-uploads sensu assets using go ship done",
		Version: version,
		RunE:    conf.Run,
	}

	flags := cmd.Flags()
	flags.StringVarP(&conf.Source, "source", "s", "", "Source Asset URL (required)")
	flags.StringVarP(&conf.Target, "target", "t", "", "Target URL stub (required)")
	flags.StringVarP(&conf.Upload, "upload", "u", "", "Upload to (scp target, required)")
	flags.StringVarP(&conf.Workdir, "workdir", "w", "", "Work dir, not deleted (default: random temporary dir, deleted)")

	return cmd, nil
}

func (conf *redistConfig) check() error {
	for _, item := range []struct {
		name string
		raw  string
		url  **url.URL
	}{
		{"--source", conf.Source, &conf.srcURL},
		{"--target", conf.Target, &conf.tarURL},
		{"--upload", conf.Upload, nil},
	} {
		if item.raw == "" {
			return fmt.Errorf("please provide %s", item.name)
		}

		if item.url == nil {
			continue
		}

		parsed, err := url.Parse(item.raw)
		if err != nil {
			return fmt.Errorf("%s: %w", item.name, err)
		}

		*item.url = parsed
	}

	return nil
}

func (conf *redistConfig) Run(cmd *cobra.Command, args []string) error {
	if err := conf.check(); err != nil {
		return err
	}

	if conf.Workdir == "" {
		dir, err := ioutil.TempDir("", "sensu-redist-asset")
		if err != nil {
			return fmt.Errorf("creating temp directory for sensu-asset: %w", err)
		}

		defer os.RemoveAll(dir)

		conf.Workdir = dir
	}

	r, err := redist.New(conf.Workdir)
	if err != nil {
		return err
	}

	if err := r.Download(conf.srcURL); err != nil {
		return err
	}

	if err := r.Sanity(); err != nil {
		return err
	}

	if err := r.RelocateTo(conf.tarURL); err != nil {
		return err
	}

	if err := r.Upload(conf.Upload); err != nil {
		return err
	}

	fmt.Printf("The asset is now available at %s\n", r.Asset.New.String())

	return nil
}
