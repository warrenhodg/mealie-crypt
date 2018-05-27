package main

import (
	"fmt"
	"github.com/warrenhodg/go_licenses"
	"gopkg.in/alecthomas/kingpin.v2"
)

func setupLicenseCommand(app *kingpin.Application) {
	app.Command("license", "Show the license")
}

func handleLicenseCommand(commands []string) error {
	fmt.Printf(license.Mit(copyrightYear, copyrightHolder))
	return nil
}
