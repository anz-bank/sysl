.PHONY: all init

all:  \
	ecosystem-test.all \
	foo.bar

init:
	mkdir -p .blueprint/children


### arrai-demo (self)

.PHONY: arrai-demo.make
arrai-demo.make: 
	cd ./ && make

### foo (local)


.PHONY: foo.bar
foo.bar: 
	cd ./out && ls

### ecosystem-test

.PHONY: ecosystem-test.clone
ecosystem-test.clone: init
	rm -rf .blueprint/children/ecosystem-test && \
		git clone --depth=1 https://github.com/orlade-anz/ecosystem-test .blueprint/children/ecosystem-test

.PHONY: ecosystem-test.all
ecosystem-test.all: ecosystem-test.clone
	cd .blueprint/children/ecosystem-test && make
