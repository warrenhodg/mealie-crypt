package main

func addPubKey(file *string, alias *string, keyFile *string, comment *string) error {
    var teamPassKey TeamPassKey

    teamPassFile, err := readFile(file)
    if err != nil {
        return err
    }

    keyValue, err := readKey(keyFile)
    if err != nil {
        return err
    }

    teamPassKey.Alias = *alias
    teamPassKey.Value = keyValue
    teamPassKey.Comment = *comment

    teamPassFile.PublicKeys = append(teamPassFile.PublicKeys, teamPassKey)

    return writeFile(file, teamPassFile)
}
