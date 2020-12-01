.PHONY: build start lint tidy

all: build

build:
	npm build

start:
	npm start

lint:
	find ./docs -name '*.md' | xargs -1 npx markdown-link-check --progress

tidy:
	yarn format
