package cmd

// containsString search for a value in a list of strings
// return a boolean true or false if found.  If found, also
// return  it's index / offset
func containsString(val string, slice []string) (bool, int) {
	for i, v := range slice {
		if val == v {
			return true, i
		}
	}
	return false, 0
}

// flattenUniqueStrings given a list of strings return
// all those that are unique and flatten into a single
// string
func flattenUniqueStrings(list []string) string {
	var uniqueItems []string
	var response string

	for _, v := range list {
		if repeated, _ := containsString(v, uniqueItems); repeated == true {
			continue
		}
		uniqueItems = append(uniqueItems, v)
	}

	for _, v := range uniqueItems {
		response += v + "\n"
	}

	return response
}
