SYSL_PLANTUML=http://localhost:8080

all: data_diagram data_dictionary integration_diagram terraform_gdc

init:
	mkdir -p out

data_diagram: init
	arrai run $@.arrai > out/$@.txt
	curl $(SYSL_PLANTUML)/svg/~1$(shell cat out/$@.txt | plantuml-encode) > out/$@.svg

data_dictionary: init
	arrai run $@.arrai > out/$@.txt

integration_diagram: init
	arrai run $@.arrai > out/$@.txt

terraform_gdc: init
	arrai run $@.arrai > out/$@.txt
