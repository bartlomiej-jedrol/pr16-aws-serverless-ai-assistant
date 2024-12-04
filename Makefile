ASSISTANT_PATH = "cmd/assistant"
BIN_PATH = "bin"

clean:
	cd ${BIN_PATH} && rm -f bootstrap

build: clean
	cd ${ASSISTANT_PATH} && GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
	mv ${ASSISTANT_PATH}/bootstrap ${BIN_PATH}
	cd ${BIN_PATH} && zip assistant.zip bootstrap

deploy: build
	aws lambda update-function-code --function-name pr16-assistant \
	--zip-file fileb://${BIN_PATH}/assistant.zip
