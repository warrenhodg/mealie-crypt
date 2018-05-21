package main

func initFile(filename *string, comment *string) error {
    var teamPassFile TeamPassFile

    teamPassFile.Comment = *comment

    return writeFile(filename, teamPassFile)
}
