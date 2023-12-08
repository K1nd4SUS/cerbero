package metricsregex

import "encoding/hex"

func ToHex(regex string) string {
	return hex.EncodeToString([]byte(regex))
}
