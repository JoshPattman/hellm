default:
	echo "Please run either make tool, make extension, or make final-build (reccomended)"
tool:
	go build -o hellm .
	chmod +x hellm
	mv hellm /usr/local/bin
extension:
	rm -rf ~/.vscode/extensions/hellm
	cp -r ./vscode/hellm ~/.vscode/extensions/hellm
	rm -rf ~/.cursor/extensions/hellm
	cp -r ./vscode/hellm ~/.cursor/extensions/hellm
final-build:
	rm -rfd bin/
	mkdir bin
	cd vscode/hellm; vsce package
	mv vscode/hellm/hellm-0.0.1.vsix bin/
	go build -o bin/hellm .
	GOOS=windows GOARCH=amd64 go build -o bin/hellm.exe .