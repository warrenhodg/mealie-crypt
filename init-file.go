package main

import (
    "gopkg.in/alecthomas/kingpin.v2"
)

func addInitFileCommand(app *kingpin.Application) {
    app.Command("init-file", "Initialize a new empty file")
}

func initFile() error {
    var teamPassFile TeamPassFile

    teamPassFile.Comment = *comment

    return writeFile(filename, true, teamPassFile)
}
