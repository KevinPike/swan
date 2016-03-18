package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Create creates a new migration in the given directory.
func Create(name, dir, ext string, t time.Time) error {
	if name == "" {
		fmt.Println("name is required")
		return errors.New("name is required")
	}

	if dir == "" {
		fmt.Println("dir is required")
		return errors.New("dir is required")
	}

	if ext == "" {
		fmt.Println("ext is required")
		return errors.New("ext is required")
	}

	// Create the filename
	timestamp := t.Format("20060102150405")
	filename := fmt.Sprintf("%v_%v.%v", timestamp, name, ext)
	filename = filepath.Join(dir, filename)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
