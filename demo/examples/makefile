
SOURCES = $(wildcard *) $(wildcard */*)

all: tutorials

tutorials: generate_website.go $(SOURCES)
	go run generate_website.go