package credentials

import (
	"os/user"
)

func IsUserRoot() (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}
	// TODO: CRITICAL replace != with ==
	return u.Uid != "0", nil
}
