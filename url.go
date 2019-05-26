package getrelease

import (
	"fmt"
	"regexp"
	"strings"
)

// resolveChecksum downloads checksum file and creates url query compatible with go-getter
func resolveChecksum(urlGetter urlGetter, checksum string) (string, error) {
	vs := strings.SplitN(checksum, ":", 2)
	switch len(vs) {
	case 2:
		break // good
	default:
		// if no checksum type us provided, it is invalid
		return "", fmt.Errorf("unsupported checksum: %s", checksum)
	}

	checksumType, checksumValue := vs[0], vs[1]

	switch checksumType {
	case "asset":
		reg, err := regexp.Compile(checksumValue)
		if err != nil {
			return "", err
		}

		url, err := urlGetter(reg)
		if err != nil {
			return "", err
		}

		return "file:" + *url, nil
	case "md5", "sha1", "sha256", "sha512":
		return checksum, nil
	case "link", "file":
		return "file:" + checksumValue, nil
	default:
		return "", fmt.Errorf("unsupported checksum type: %s", checksumType)
	}
}

// adjustUrlForGetter creates url queries compatible with go-getter
func adjustUrlForGetter(url string, urlGetter urlGetter, option *Option) (string, error) {
	url += "?"

	if strings.TrimSpace(option.Archive) != "" {
		url += fmt.Sprintf("%s=%s&", "archive", option.Archive)
	}
	if strings.TrimSpace(option.Checksum) != "" {
		checksum, err := resolveChecksum(urlGetter, option.Checksum)
		if err != nil {
			return "", err
		}
		url += fmt.Sprintf("%s=%s&", "checksum", checksum)
	}
	if strings.TrimSpace(option.FileName) != "" {
		url += fmt.Sprintf("%s=%s&", "filename", option.FileName)
	}

	return strings.Trim(strings.Trim(url, "?"), "&"), nil
}
