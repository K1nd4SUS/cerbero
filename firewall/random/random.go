package random

import "math/rand"

const (
	LowerAlphabet   = "abcdefghijklmnopqrstuvwxyz"
	UpperAlphabet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DigitAlphabet   = "1234567890"
	SpecialAlphabet = "!@#$%^&*"
)

func String(alphabet string, length int) string {
	res := ""
	for i := 0; i < length; i++ {
		res += string(alphabet[rand.Intn(len(alphabet))])
	}
	return res
}
