mealie-crypt: windows linux mac
windows: mealie-crypt.exe
linux: mealie-crypt.linux
mac: mealie-crypt.mac

mealie-crypt.exe: *.go
	GOOS=windows GOARCH=386 go build -o mealie-crypt.exe .
mealie-crypt.linux: *.go
	GOOS=linux GOARCH=386 go build -o mealie-crypt.linux .
mealie-crypt.mac: *.go
	GOOS=darwin GOARCH=386 go build -o mealie-crypt.mac .

clean:
	rm -f mealie-crypt.exe mealie-crypt.linux mealie-crypt.mac
