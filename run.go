package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Run runs swan migrations.
func Run(filename, dir string) error {
	// Open the last migration file
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("unable to read or create file", filename)
		return err
	}
	defer f.Close()

	// Read the last migration stored in it.
	var lastB []byte
	lastB, err = ioutil.ReadAll(f)
	if err != nil {

	}
	last := string(lastB)
	last = strings.Trim(last, "\n\r")

	// skipFirstMigration is used to indicate skipping the first migration. This
	// is used because if the last migration is x, we want to run everything after x,
	// not including x. However, if there was no last migration we want to run all
	// migrations and not skip the first one.
	skipFirstMigration := last != ""

	// Read migrations from directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {

	}

	// Extract file names from file infos.
	migrations := make([]string, len(files))
	for i, file := range files {
		migrations[i] = file.Name()
	}

	// Sort migrations so they can be run in order
	sort.Strings(migrations)

	// Execute later migrations
	i := sort.SearchStrings(migrations, last)
	if i > len(migrations) {
		fmt.Println("No such migration", last)
		return errors.New("no such migration")
	}

	if skipFirstMigration {
		i++
	}

	// Are there no migrations to run?
	if i == len(migrations) {
		fmt.Println("up to date")
		return nil
	}

	// Execute each migration
	for ; i < len(migrations); i++ {
		migration := "migrations/" + migrations[i]
		if err := exec.Command(migration).Run(); err != nil {
			fmt.Println("error: ", migration, err)
			break
		}
		fmt.Println("executed", migration)
	}

	// if i is 0, we haven't run a migration before and we failed to run the first
	// migration.
	if i == 0 {
		return errors.New("no migration executed")
	}

	// Empty the file
	if err := f.Truncate(0); err != nil {
		fmt.Println("unable to truncate file", err)
		return err
	}

	// Seek to the beginning of the file
	if _, err := f.Seek(0, 0); err != nil {
		fmt.Println("unable to seek file", err)
		return err
	}

	// Write the migration to the last migration file
	lastMigration := migrations[i-1]
	if _, err := f.WriteString(lastMigration); err != nil {
		fmt.Println("unable to write last migration", lastMigration, "to", filename)
		return err
	}
	return nil
}
