# Mealie-crypt
A team, and public repository safe method for storing and managing sensitive information.

## Usage
Typical flow of the application is:
- Create an empty file:
  - Quickly create `mealie-crypt.yaml`: `mealie-crypt file init`
  - Create some other file : `mealie-crypt file init -f my-file.yaml`
    
- Adding users:
  - Add myself as a user to the file: `mealie-crypt users add`
  - Add some other user to the file: `mealie-crypt users add -u their-name -K their-pub-key.pub`

- Create groups:
  - Create group `_`, and add myself as a user: `mealie-crypt groups add`
  - Create some other group: `mealie-crypt groups add -g my-name`
  - Create group `_`, with specific users: `mealie-crypt groups add -U user1 -U user2`

- Adding users to groups:
  - Add user to existing group: `mealie-crypt group user-add -u your-name -U their-name`
    - *You must be a part of the group to which you are adding users.*

- Add values:
  - Add a value to group `_`: `mealie-crypt values set -n foo -v bar`
  - Decrypt, edit and re-encrypt:
    - `mealie-crypt decrypt`
    - edit the `decrypted` object in the file
    - `mealie-crypt encrypt`

- Search for stuff: `mealie-crypt values get -n '*stuff*'`

## Security
Mealie-crypt works by:
- Storing one or more user's public RSA keys `PUB_KEY` (recommended is to use 2048 bit encrypted private keys)
- Creating a 256 bit key `SYM_KEY` per group
- Encrypting the group's key `SYM_KEY` with each user's public key `PUB_KEY` using OAEP algorithm, and storing those with the group
- Encrypting each value in the group with the symmetrical key using AES-256 `SYM_KEY`

## Repository usage
To use the encrypted file in a team, the following goals have been met:
- Users do not need to share their passwords - only their private keys
- The structure of the file is text-based yaml, which handles well in git merge functions, and is human-editable too
- It is possible to mass decrypt, and encrypt the file - preserving encrypted content that does not change, so as to minimize changed-line counts.

