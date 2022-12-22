package services

import "strings"

func ParseHashtags(content string) (tags []string) {
	tags = make([]string, 0)
	words := strings.Fields(content)

	// Iterate over the words and check for hashtags
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#") // trim '#'
			tags = append(tags, strings.ToLower(tag))
		}
	}

	return
}
