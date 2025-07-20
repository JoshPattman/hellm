default:
	echo "Please run either make tool or make extension"
tool:
	go build -o hellm .
	chmod +x hellm
	mv hellm /usr/local/bin
extension:
	rm -rf ~/.vscode/extensions/hellm
	cp -r ./vscode/hellm ~/.vscode/extensions/hellm
	rm -rf ~/.cursor/extensions/hellm
	cp -r ./vscode/hellm ~/.cursor/extensions/hellm