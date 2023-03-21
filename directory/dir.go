package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

func Generate(parentPath string, sub any) error {
	if parentPath != "." {
		if err := os.MkdirAll(parentPath, 0755); err != nil {
			return err
		}
		defaultFilePath := filepath.Join(parentPath, filepath.Base(parentPath)+".go")
		if _, err := os.Stat(defaultFilePath); os.IsNotExist(err) {
			f, err := os.Create(defaultFilePath)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(f, "package %s", filepath.Base(parentPath)); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	switch v := sub.(type) {
	case map[string]any:
		for key, value := range v {
			err := Generate(filepath.Join(parentPath, key), value)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}
