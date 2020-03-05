package template

type Vars struct {
	Service Service `json:"service"`
}

type ServiceMap map[string]Service

type serviceBase struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Service struct {
	serviceBase
	Port string      `json:"port"`
	Impl ServiceImpl `json:"impl"`
}

type ServiceImpl struct {
	serviceBase
	InterfaceFactory string `json:"interface_factory"`
	CallbackFactory  string `json:"callback_factory"`
}
