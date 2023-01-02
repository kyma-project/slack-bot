IMG ?= eu.gcr.io/kyma-project/slack-bot

.PHONY: build-image
build-image:
	docker build . -t ${IMG}

.PHONY: build-image
push-image:
	docker push ${IMG}

.PHONY: build-cleaner
build-cleaner:
	go build -o bot main.go
