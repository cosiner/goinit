package goinit

func removeBrackets(s string) string {
	l := len(s)
	if l >= 2 && (s[0] == '(' && s[l-1] == ')') {
		return s[1 : l-1]
	}
	return s
}