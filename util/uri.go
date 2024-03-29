package util

import (
	"errors"
	"strings"
)

// -----------------------------------------------------------------------------

func SanitizeUrlPath(path string) (string, error) {
	// Nothing to sanitize?
	if len(path) == 0 {
		return "/", nil
	}

	// Convert backslashes
	path = strings.ReplaceAll(path, "\\", "/")

	// Check if a trailing slash is present
	endsWithSlash := strings.HasSuffix(path, "/")

	// Split path in fragments and remove empty
	newPathFragments := make([]string, 0)
	pathFragments := strings.Split(path, "/")
	for _, frag := range pathFragments {
		if len(frag) > 0 && frag != "." {
			if frag != ".." {
				newPathFragments = append(newPathFragments, frag)
			} else {
				if len(newPathFragments) > 0 {
					newPathFragments = newPathFragments[0 : len(newPathFragments)-1]
				} else {
					return "", errors.New("invalid path")
				}
			}
		}
	}

	// Join again
	path = "/" + strings.Join(newPathFragments, "/")
	if len(newPathFragments) > 0 && endsWithSlash {
		path += "/"
	}

	// Done
	return path, nil
}
