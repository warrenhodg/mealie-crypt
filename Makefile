teampass: windows linux mac
windows: teampass.exe
linux: teampass.linux
mac: teampass.mac

teampass.exe: *.go
	GOOS=windows GOARCH=386 go build -o teampass.exe .
teampass.linux: *.go
	GOOS=linux GOARCH=386 go build -o teampass.linux .
teampass.mac: *.go
	GOOS=darwin GOARCH=386 go build -o teampass.mac .

clean:
	rm -f teampass.exe teampass.linux teampass.mac
