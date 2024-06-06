package did

import (
	"fmt"
	"strings"
)

func IsValidDID(did string) error {
	_, _, _, err := parseDID(did)
	if err != nil {
		return err
	}
	return nil
}

func parseDID(did string) (scheme string, method string, id string, err error) {
	/*
		const idchar = `a-zA-Z0-9-_\.`
		regex := fmt.Sprintf(`^did:[a-z0-9]+:(:+|[:%s]+)*[%%:%s]+[^:]$`, idchar, idchar)

		r, err := regexp.Compile(regex)
		if err != nil {
			panic("err")
		}

		if !r.MatchString(did) {
			err = fmt.Errorf("invalid did: %s", did)
			return
		}
	*/

	parts := strings.SplitN(did, ":", 3)
	if len(parts) != 3 {
		err = fmt.Errorf("%v is not a valid did", did)
		return
	}

	scheme = "did"
	method = parts[1]
	id = parts[2]
	return
}
