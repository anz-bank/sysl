all: diagrams scripts

.PHONY: diagrams scripts

diagrams: data ints epa sd
scripts: db

data:
	sysl data --output out/data.png --project "Bank :: Data Views" ./project.sysl

ints:
	sysl ints --output out/ints.png --project "Bank :: Integrations" ./project.sysl

epa:
	sysl ints --epa --output out/epa.png --project "Bank :: Integrations" ./project.sysl

sd:
	sysl sd --output out/sd.png --endpoint "ATM <- Withdraw" ./project.sysl

db:
	sysl generate-db-scripts --output-dir out --app-names BankModel ./bank.sysl
