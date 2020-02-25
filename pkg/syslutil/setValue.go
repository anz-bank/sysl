package syslutil

func ResetVal(init string, new string) string {
	if init == "" {
		return new
	}

	return init
}
