package redist

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

func (r *Redist) Sanity() error {
	for _, build := range r.Builds {
		f, err := os.Open(build.Local)
		if err != nil {
			return fmt.Errorf("opening %q for checksumming: %w", build.Local, err)
		}

		hasher := sha512.New()

		if _, err := io.Copy(hasher, f); err != nil {
			return fmt.Errorf("reading %q for checksumming: %w", build.Local, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("closing %q after checksumming: %w", build.Local, err)
		}

		computed := fmt.Sprintf("%x", hasher.Sum(nil))
		if computed != build.SHA512 {
			return fmt.Errorf("checksum error for %q: got %s, want %s", build.Local, computed, build.SHA512)
		}
	}

	return nil
}
