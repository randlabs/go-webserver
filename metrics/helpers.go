package metrics

// -----------------------------------------------------------------------------

func isJSON(s string) bool {
	// An official (?) method but a plain text is also considered a valid JSON
	// var js interface{}
	// return json.Unmarshal([]byte(s), &js) == nil

	// Our quick approach
	startIdx := 0
	endIdx := len(s)

	// Skip blanks at the beginning and the end
	for startIdx < endIdx && isBlank(s[startIdx]) {
		startIdx += 1
	}
	for endIdx > startIdx && isBlank(s[endIdx-1]) {
		endIdx -= 1
	}

	return startIdx < endIdx &&
		((s[startIdx] == '{' && s[endIdx-1] == '}') ||
			(s[startIdx] == '[' && s[endIdx-1] == ']'))
}

func isBlank(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}
