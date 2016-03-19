package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/mitchellh/iochan"
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Kill, os.Interrupt)

	// Execute each migration
	var quit bool
	for ; i < len(migrations); i++ {
		quit = false
		outCh := make(chan string, 1)
		doneCh := make(chan struct{}, 1)
		go func() {
			for line := range outCh {
				fmt.Println(line)
			}
			close(doneCh)
		}()

		migration := filepath.Join(dir, migrations[i])
		outCh <- "==> executing " + migration
		cmd := exec.Command(migration)

		if err := execute(cmd, outCh, sigCh); err != nil {
			outCh <- "==> " + err.Error()
			quit = true
		}

		// Signal that there are no more lines to print
		close(outCh)
		// Wait for lines to finish printing
		<-doneCh
		if quit {
			break
		}
		fmt.Println("==> OK\t" + migration)
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

func execute(cmd *exec.Cmd, outCh chan<- string, sigCh <-chan os.Signal) error {
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to read stdout from command: %s", err)
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to read stderr from command: %s", err)
	}

	// Create the channels we'll use for data
	exitCh := make(chan int, 1)
	doneCh := make(chan interface{}, 1)
	stdoutCh := iochan.DelimReader(outPipe, '\n')
	stderrCh := iochan.DelimReader(errPipe, '\n')
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to start command: %s", err)
	}

	go func() {
		select {
		case <-doneCh:
			return
		case sig := <-sigCh:
			cmd.Process.Signal(sig)
		}
	}()

	// Start the goroutine to watch for the exit
	go func() {
		exitStatus := 0

		err := cmd.Wait()
		doneCh <- struct{}{}

		if exitErr, ok := err.(*exec.ExitError); ok {
			exitStatus = 1

			// There is no process-independent way to get the REAL
			// exit status so we just try to go deeper.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitStatus = status.ExitStatus()
			}
		}

		exitCh <- exitStatus
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	streamFunc := func(ch <-chan string) {
		defer wg.Done()
		for data := range ch {
			if data != "" {
				outCh <- data
			}
		}
	}

	go streamFunc(stdoutCh)
	go streamFunc(stderrCh)

	exitStatus := <-exitCh
	wg.Wait()

	if exitStatus != 0 {
		return fmt.Errorf("non-zero exit code: %d", exitStatus)
	}
	return nil
}
