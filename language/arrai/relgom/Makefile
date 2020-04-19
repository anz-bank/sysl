all:
	sysl pb --mode=json ../../../demo/petshop/petshop.sysl | \
	./relgom.arrai PetShopModel | \
	tar xvf - -C out && goimports -w ./out && colordiff -u ../../go/cmd/relgom/example/petshop/ out/
