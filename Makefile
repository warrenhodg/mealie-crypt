dioscorea: windows linux mac
windows: dioscorea.exe
linux: dioscorea.linux
mac: dioscorea.mac

dioscorea.exe: *.go
	GOOS=windows GOARCH=386 go build -o dioscorea.exe .
dioscorea.linux: *.go
	GOOS=linux GOARCH=386 go build -o dioscorea.linux .
dioscorea.mac: *.go
	GOOS=darwin GOARCH=386 go build -o dioscorea.mac .

clean:
	rm -f dioscorea.exe dioscorea.linux dioscorea.mac
