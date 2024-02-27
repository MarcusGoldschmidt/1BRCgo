package unit

import "strconv"

type Size int64

const (
	B  Size = 1
	KB      = B * 1024
	MB      = KB * 1024
	GB      = MB * 1024
	TB      = GB * 1024
	PB      = TB * 1024
)

func (v Size) ToString() string {
	switch {
	case v < KB:
		return strconv.FormatInt(int64(v), 10) + "B"
	case v < MB:
		return strconv.FormatInt(int64(v/KB), 10) + "KB"
	case v < GB:
		return strconv.FormatInt(int64(v/MB), 10) + "MB"
	case v < TB:
		return strconv.FormatInt(int64(v/GB), 10) + "GB"
	case v < PB:
		return strconv.FormatInt(int64(v/TB), 10) + "TB"
	default:
		return strconv.FormatInt(int64(v/TB), 10) + "TB"
	}
}
