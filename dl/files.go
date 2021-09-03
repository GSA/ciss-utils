package dl

import (
	"fmt"
	"io"
	"os"
)

func save(r io.Reader, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("failed to close %s file handle: %v", path, err)
		}
	}()
	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("failed to save inventory: %v", err)
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)

	// return the negative of is not exist
	return !os.IsNotExist(err)
}
