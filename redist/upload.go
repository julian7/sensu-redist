package redist

import "os/exec"

func (r *Redist) Upload(target string) error {
	args := make([]string, 0, len(r.Builds)+2)

	args = append(args, r.Asset.Local)

	for _, build := range r.Builds {
		args = append(args, build.Local)
	}

	return exec.Command("scp", append(args, target)...).Run()
}
