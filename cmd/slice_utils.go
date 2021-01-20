package cmd

func containsString(val string, slice []string) (bool, int) {
	for i, v := range slice {
		if val == v {
			return true, i
		}
	}
	return false, 0
}
