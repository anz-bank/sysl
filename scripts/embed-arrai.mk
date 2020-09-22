embed-arrai:
	cd ./pkg/importer/avro && ./avro.arrai.go.sh \
	&& cd ./../spanner && ./spanner.arrai.go.sh
