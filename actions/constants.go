package actions

const (
	yamlFileRegexStr       = `(?mi).*\.(yaml|yml)`
	dockerImageRefRegexStr = `(?m)image:\s*(?P<image>[^[{\s]+)\s+`
)
