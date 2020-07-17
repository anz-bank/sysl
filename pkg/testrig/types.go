package testrig

type ServiceVarsMap map[string]ServiceVars

type serviceVarsBase struct {
	Name   string `json:"name"`
	Import string `json:"import"`
	URL    string `json:"url"`
}

type ServiceVars struct {
	serviceVarsBase
	Port string          `json:"port"`
	Impl ServiceImplVars `json:"impl"`
}

type ServiceImplVars struct {
	serviceVarsBase
	InterfaceFactory string `json:"interface_factory"`
	CallbackFactory  string `json:"callback_factory"`
}
