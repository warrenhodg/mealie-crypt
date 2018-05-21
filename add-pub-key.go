package main

func addPubKey(file *string, alias *string, keyFile *string, comment *string) error {
    var teamPassKey TeamPassKey

    teamPassFile, err := readFile(file)
    if err != nil {
        return err
    }

    keyContent, err := readKey(keyFile)
    if err != nil {
        return err
    }

    teamPassKey.Alias = *alias
    teamPassKey.Key = keyContent
    teamPassKey.Comment = *comment

    teamPassFile.Keys = append(teamPassFile.Keys, teamPassKey)

    return writeFile(file, teamPassFile)
}
