package utility

import (
	"fmt"
	"errors"
)

const (
	DEV = "dev"
	TEST = "test"
	PROD = "prod"
)

var ENVs = []string{DEV, TEST, PROD}

func Verify(env string, def string) (string, error) {
	if env == "" {
		return def, nil
	}

	if env != DEV && env != TEST && env != PROD {
		return "", errors.New(fmt.Sprintf("环境只能是%s其中之一,默认为" + DEV + ".", ENVs))
	}

	return env, nil
}