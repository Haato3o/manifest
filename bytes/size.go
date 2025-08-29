package bytes

import (
	"errors"
	"regexp"
	"strconv"
)

const (
	Byte     = 1
	Kilobyte = 1 << 10
	Megabyte = 1 << 20
	Gigabyte = 1 << 30
)

var (
	sizeRegex = regexp.MustCompile(`(?i)^(\d+(?:\.\d+)?)\s*(B|KB|MB|GB)$`)
)

type Size struct {
	bytes int
}

func (s Size) Bytes() int {
	return s.bytes
}

func ParseSize(s string) (Size, error) {
	if !sizeRegex.MatchString(s) {
		return Size{}, errors.New("invalid size format")
	}

	matches := sizeRegex.FindStringSubmatch(s)
	sizeRaw := matches[1]
	byteFormat := matches[2]

	size, err := strconv.ParseFloat(sizeRaw, 64)
	if err != nil {
		return Size{}, err
	}

	switch byteFormat {
	case "B":
		size = size * Byte
	case "KB":
		size = size * Kilobyte
	case "MB":
		size = size * Megabyte
	case "GB":
		size = size * Gigabyte
	default:
		break
	}

	return Size{int(size)}, nil
}
