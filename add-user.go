package main

func addUser(file *string, name *string, keyFile *string, comment *string) error {
    var user TeamPassUser

    teamPassFile, err := readFile(file, true)
    if err != nil {
        return err
    }

    keyValue, err := readKey(keyFile)
    if err != nil {
        return err
    }

    user.Name = *name
    user.Value = keyValue
    user.Comment = *comment

    teamPassFile.Users = append(teamPassFile.Users, user)

    return writeFile(file, false, teamPassFile)
}
