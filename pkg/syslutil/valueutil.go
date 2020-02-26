package syslutil

func GetNonEmpty(init string, new string) string {
	if init == "" {
		return new
	}

	return init
}
