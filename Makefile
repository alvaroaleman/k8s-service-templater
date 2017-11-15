SHELL = /usr/bin/env bash

k8s-service-templater: main.go
	@docker run --rm \
			-v $$PWD:/go/src/k8s-service-templater \
			-w /go/src/k8s-service-templater \
			golang:1.8 \
			env CGO_ENABLED=0 go build -v

app-container: k8s-service-templater
	@docker run -d \
		-it \
		-v $$PWD:/usr/local/bin:ro \
		-w /usr/local/bin \
		-p 127.0.0.1:5000:5000 \
		--name=k8s-service-templater-app-container \
		alpine:latest \
		/usr/local/bin/k8s-service-templater &>/dev/null

bash_unit:
	curl -LO https://raw.githubusercontent.com/pgrange/bash_unit/v1.6.0/bash_unit
	chmod +x bash_unit

.PHONY: test
test: bash_unit
	./bash_unit test/integration.sh


clean:
	- @docker rm -f k8s-service-templater-app-container &>/dev/null
