.PHONY: install build start lint tidy update

all: build

install:
	yarn install --frozen-lockfile

build:
	yarn build

start:
	yarn start

lint:
	find ./{blog,docs} -name '*.md' | xargs -n1 npx markdown-link-check --progress --config=.md-link-check-config.json

tidy:
	yarn format

update:
	npx npm-check-updates -u && yarn install
