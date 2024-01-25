package strconv

import (
	"strconv"
)

const (
	// DecimalBase used for string parsing.
	DecimalBase = 10
	// Int64BitSize used for string parsing.
	Int64BitSize = 64
)

func ParseDecimalInt(s string) (int64, error) {
	return strconv.ParseInt(s, DecimalBase, Int64BitSize)
}

func FormatDecimalInt(i int64) string {
	return strconv.FormatInt(i, DecimalBase)
}
