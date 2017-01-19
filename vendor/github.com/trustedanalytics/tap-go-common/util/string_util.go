package util

import (
	"fmt"
	"regexp"
)

var imageAdressRegexp = regexp.MustCompile(`(?P<host>[^:]+:?.+)/(?P<name>[^:]+):?(?P<tag>.*)$`)

func ParseImageAddress(imageAddress string) (hostWithPort, imageName, imageTag string, err error) {
	match := imageAdressRegexp.FindStringSubmatch(imageAddress)
	result := make(map[string]string)
	if len(match) == 4 {
		for i, name := range imageAdressRegexp.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}
		hostWithPort = result["host"]
		imageName = result["name"]
		imageTag = result["tag"]
	} else {
		err = fmt.Errorf("cannot split image value from address: %s using regexp: %s", imageAddress, imageAdressRegexp)
	}
	return
}
