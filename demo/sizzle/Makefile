basepath=github.com/sysl/demo/sizzle/server
SYSL_GO_ROOT = $(HOME)/github.com/anz-bank/sysl-go

all: $(foreach s,Identity Product Order,servers/$(s)Server)

servers/%: sizzle.json
	arrai --out=dir:$@ $(SYSL_GO_ROOT)/codegen/arrai/auto.arrai $(basepath) $< $* =

sizzle.json: sizzle.sysl
	sysl pb --mode=json $< > $@
