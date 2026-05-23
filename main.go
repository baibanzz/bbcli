package main

import (
	"os"

	"github.com/baibanzz/bbcli/internal/windows"
)

func main() {
	path, _ := os.Getwd()
	windows.NewApp(path).Run()
}