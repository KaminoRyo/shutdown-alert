package config

// TargetURL is the URL to be opened.
const TargetURL = "https://www.google.com" // As specified in the requirements document

// GetTargetURL returns the hardcoded target URL.
// This is a pure function.
func GetTargetURL() string {
	return TargetURL
}
