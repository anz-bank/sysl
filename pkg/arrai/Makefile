.PHONY: all clean init init-diagram plantuml data_diagram data_dictionary integration_diagram terraform_gdc

SYSL_PLANTUML=http://localhost:8080
OUT=gen

all: data_diagram data_dictionary integration_diagram terraform_gdc

# Deletes all generated outputs.
clean:
	rm -f $(OUT)/*

# Metadata

data_dictionary: init
	arrai run $@.arrai > $(OUT)/$@.csv

terraform_gdc: init
	arrai run $@.arrai > $(OUT)/$@.tf

# Diagrams

data_diagram: init-diagram
	arrai run $@.arrai > $(OUT)/$@.puml && \
	curl $(SYSL_PLANTUML)/svg/~1$(shell cat $(OUT)/$@.puml | plantuml-encode) > $(OUT)/$@.svg

integration_diagram: init-diagram
	arrai run $@.arrai > $(OUT)/$@.puml && \
	curl $(SYSL_PLANTUML)/svg/~1$(shell cat $(OUT)/$@.puml | plantuml-encode) > $(OUT)/$@.svg

# Platform

# Initializes the environment to support other targets.
init:
	mkdir -p $(OUT)

# Initializes the environment to support diagram generation targets.
init-diagram: init plantuml
	which plantuml-encode || go install github.com/orlade/plantuml-encode

plantuml:
	# Version pinned to reduce variability of generated outputs.
	docker start plantuml || docker run -d -p 8080:8080 --name=plantuml plantuml/plantuml-server:jetty-v1.2020.14
