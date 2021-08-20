module github.com/anz-bank/sysl

go 1.16

replace github.com/spf13/afero => github.com/anz-bank/afero v1.2.4

require (
	aqwari.net/xml v0.0.0-20210331023308-d9421b293817
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20210521184019-c5ad59b459ec
	github.com/anz-bank/go-bindata v3.22.0+incompatible
	github.com/anz-bank/golden-retriever v0.5.0
	github.com/anz-bank/mermaid-go v0.1.1
	github.com/anz-bank/pkg v0.0.38
	github.com/arr-ai/arrai v0.283.0
	github.com/arr-ai/frozen v0.19.0
	github.com/arr-ai/proto v0.0.0-20180422074755-2ffbedebee50
	github.com/arr-ai/wbnf v0.34.0
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cornelk/hashmap v1.0.1
	github.com/getkin/kin-openapi v0.8.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/swag v0.19.15
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.12
	github.com/kevinburke/go-bindata v3.22.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/rjeczalik/notify v0.9.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.4.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
