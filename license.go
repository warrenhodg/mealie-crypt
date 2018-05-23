package main

import (
    "fmt"
    "gopkg.in/alecthomas/kingpin.v2"
    "github.com/warrenhodg/go_licenses"
)

func addLicenseCommand(app *kingpin.Application) {
    app.Command("license", "Show the license")
}

func showLicense() error {
    fmt.Printf(license.Mit(copyrightYear, copyrightHolder))
    return nil
}
