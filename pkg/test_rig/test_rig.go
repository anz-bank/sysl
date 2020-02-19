package test_rig

type TestRig struct {
	mainContent string
	dockerBlock string
}

type TestRigGenerator interface {
	generateMain(template string, vars map[string]interface{}) string
	generateDocker(template string, vars map[string]interface{}) string
}
