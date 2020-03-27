all: sysl

input = api/sysl/simple.sysl
app = simple
dependencies = jsonplaceholder # this can be a list separated by a space or left empty
outdir = gen
# Current go import path
basepath = github.com/anz-bank/sysl/demo/examples/Code-Generation

####################################################################
#                                                                  #
#                                                                  #
#                                                                  #
# START SYSL MAKEFILE: you shouldn't need to edit anything below   #
#                                                                  #
#                                                                  #
#                                                                  #
####################################################################


TRANSLOCATION = github.com/anz-bank/sysl-go/codegen/transforms
TRANSFORMS= svc_error_types.sysl svc_handler.sysl svc_interface.sysl svc_router.sysl svc_types.sysl
DOWNSTREAMTRANSFORMS = svc_client.sysl svc_error_types.sysl svc_types.sysl
GRAMMAR=github.com/anz-bank/sysl-go/codegen/grammars/go.gen.g
START=goFile
ROOT=$(shell pwd)


# Always execute these with just `make`
.PHONY: setup clean gen
sysl: clean setup gen downstream format

# Create the output directories
setup:
	mkdir -p ${outdir}/${app}
	$(foreach path, $(dependencies), $(shell mkdir -p ${outdir}/$(path)))
    $(foreach path, $(app), $(shell mkdir -p ${outdir}/$(path)))

# Generate code for the server endpoints
gen:
	$(foreach file, $(TRANSFORMS), $(shell sysl codegen --root=$(ROOT) --basepath=$(basepath)/${outdir}/ --transform $(TRANSLOCATION)/$(file) --grammar ${GRAMMAR} --start ${START} --outdir=${outdir}/${app} --app-name ${app} $(input)))

# Generate code for the downstream clients
downstream:
	$(foreach file, $(DOWNSTREAMTRANSFORMS), $(foreach downstream, $(dependencies), $(shell sysl codegen --root=$(ROOT) --basepath=$(basepath)/${outdir}/ --transform $(TRANSLOCATION)/$(file) --grammar ${GRAMMAR} --start ${START} --outdir=${outdir}/${downstream} --app-name ${downstream} $(input))))

# Formats the generated code
format:
	gofmt -s -w ${outdir}/${app}/*
	goimports -w ${outdir}/${app}/*
	$(foreach path, $(dependencies), $(shell gofmt -s -w ${outdir}/${path}/*))
	$(foreach path, $(dependencies), $(shell goimports -w ${outdir}/${path}/*))

# Remove the generated files
clean:
	rm -rf $(outdir)

# Builds the binary
server:
	go build -o bin/sysl-template main.go
