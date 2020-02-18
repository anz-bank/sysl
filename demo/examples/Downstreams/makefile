# User guide:
# make input=<inputfile> app=<appname>
# This will create a directory `app` and generate code into it
# if this script doesn't run, try running `pip2 uninstall sysl` or
# `pip uninstall sysl` as make might be running an old sysl version

TMP = .tmp# Cache the server lib directory in tmp
SERVERLIB = /var/tmp
TRANSLOCATION = .tmp/server-lib/codegen/transforms
TRANSFORMS= svc_error_types.sysl svc_handler.sysl svc_interface.sysl svc_router.sysl svc_types.sysl svc_client.sysl
DOWNSTREAMTRANSFORMS = svc_client.sysl svc_error_types.sysl svc_types.sysl

GRAMMAR=$(wildcard .tmp/server-lib/codegen/grammars/go.gen.g)
START=goFile
app=Simple



# Always execute these with just `make`
.PHONY: setup clean gen
all: setup gen downstream format tmp

# try to clone, then try to fetch and pull
setup:
	# Syncing server-lib to $(SERVERLIB)
	git clone https://github.service.anz/sysl/server-lib/ $(SERVERLIB)/server-lib || true  # Don't fail 
	cd  $(SERVERLIB)/server-lib && git fetch && git checkout tags/v0.1.9 || true
	mkdir -p $(TMP)/server-lib/
	mkdir -p gen/${app}
	# Copying server-lib to $(TMP)
	cp -rf $(SERVERLIB)/server-lib $(TMP)/
	

# Generate files with internal git service
gen:
	$(foreach file, $(TRANSFORMS), $(shell sysl codegen --dep-path github.service.anz/sysl/syslbyexample/_examples --transform $(TRANSLOCATION)/$(file) --grammar ${GRAMMAR} --start ${START} --outdir=gen/${app} --app-name ${app} $(input)))

downstream:
	$(foreach file, $(DOWNSTREAMTRANSFORMS), $(shell sysl codegen --transform $(TRANSLOCATION)/$(file) --grammar ${GRAMMAR} --start ${START} --outdir=gen/${down} --app-name ${down} $(input)))

format:
	gofmt -s -w gen/${app}/*
	goimports -w gen/${app}/*

# Remove the tmp directory after 
tmp:
	rm -rf $(TMP)

# Remove the generated files
clean:
	rm -rf ${app}
	