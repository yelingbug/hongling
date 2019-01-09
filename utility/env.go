package utility

import (
	"fmt"
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
)

const (
	DEV  = "dev"
	TEST = "test"
	PROD = "prod"
)

var ENVs = []string{DEV, TEST, PROD}

func Verify(env string, def string) (string, error) {
	if env == "" {
		return def, nil
	}

	if env != DEV && env != TEST && env != PROD {
		return "", errors.New(fmt.Sprintf("环境只能是%s其中之一,默认为"+DEV+".", ENVs))
	}

	return env, nil
}

var CacheDir = GetCacheDir()

func GetCacheDir() string {
	userdir, err := homedir.Dir()
	if err != nil {
		userdir = "." + string(os.PathSeparator)
	}

	return fmt.Sprintf("%s%s.utility%s", userdir, string(os.PathSeparator), string(os.PathSeparator))
}

func GetUserDir() string {
	if userdir, err := homedir.Dir(); err != nil {
		return "." + string(os.PathSeparator)
	} else {
		return userdir + string(os.PathSeparator)
	}
}
