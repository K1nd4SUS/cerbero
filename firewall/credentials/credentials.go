package credentials

import (
	"os/user"
)

func IsUserRoot() (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}

	return u.Uid == "0", nil
}
