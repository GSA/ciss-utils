package dl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Requires gh cli tool installed
func CreateReleasePath(outfile string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	dir := filepath.Join(wd, "release")

	err = os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("failed to create release directory: %v", err)
	}
	return filepath.Join(dir, outfile), nil
}

func DownloadDeps(deps map[string]string, depType string) error {
	for key, version := range deps {
		var err error
		if depType == "direct" {
			path, err := CreateReleasePath(filepath.Base(key))
			if err != nil {
				return err
			}
			err = Download(fmt.Sprintf(key, version), path)
		} else if depType == "gh" {
			err = DownloadGHRelease(key, version)
		} else {
			return errors.New(fmt.Sprintf("Invalid filetype provide %s", depType))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func Download(uri string, path string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("failed to parse url: %v", err)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to download zip: %v", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("failed to close HTTP response body: %v", err)
		}
	}()
	return save(resp.Body, path)
}

func DownloadGHRelease(target string, tag string) error {
	s := strings.Split(target, "#")
	filename := s[1]
	ownerRepo := s[0]
	path, err := CreateReleasePath(filename)
	if err != nil {
		return err
	}
	if FileExists(path) {
		os.Remove(path)
	}
	return Run(map[string][]string{
		"gh": {"release", "download", "-D" + filepath.Dir(path), "-R" + ownerRepo, "-p" + filename, tag},
	})
}
