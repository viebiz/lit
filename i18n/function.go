package i18n

import (
	"os"
	"path/filepath"
)

func validateSourcePath(path string) error {
	path = filepath.Clean(path)
	if _, err := os.Stat(path); err != nil {
		return err
	}

	return nil
}
