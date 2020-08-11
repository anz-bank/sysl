all:  \
	_foo.bar

init:
	mkdir -p ./blueprint/children

# arrai-demo

_arrai-demo.make:
	cd pkg/arrai && make

# foo

_foo.clone: init
    git clone https://github.com/anz-bank/sysl ./blueprint/children

_foo.bar:
	cd pkg/arrai/out && ls
