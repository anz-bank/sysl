.PHONY: all init

all:  \
	_ecosystem-test.all \
	_foo.bar

init:
	mkdir -p .blueprint/children

# arrai-demo (self)

.PHONY: _arrai-demo.make
_arrai-demo.make: 
	cd pkg/arrai && make

# foo (local)


.PHONY: _foo.bar
_foo.bar: 
	cd pkg/arrai/out && ls

# ecosystem-test 

.PHONY: _ecosystem-test.clone
_ecosystem-test.clone: init
	rm -rf .blueprint/children/ecosystem-test && git clone --depth=1 https://github.com/orlade-anz/ecosystem-test .blueprint/children/ecosystem-test

.PHONY: _ecosystem-test.all
_ecosystem-test.all: _ecosystem-test.clone
	cd .blueprint/children/ecosystem-test && make
