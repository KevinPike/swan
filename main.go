package main

import "os"

func main() {
	if err := Run(".swan", "./migrations"); err != nil {
		os.Exit(1)
	}
}
