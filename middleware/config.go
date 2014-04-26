package middleware

import (
	"errors"
	"fmt"
	"strings"
)

type MiddlewareConfig map[string]string

func (cfg MiddlewareConfig) CheckForMandatoryKeys(keys ...string) error {
	missing := make([]string, 0, len(keys))
	for _, key := range keys {
		_, ok := cfg[key]
		if !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return missingKeysError(missing)
	}

	return nil
}

func missingKeysError(keys []string) error {
	list := strings.Join(keys, ", ")
	return errors.New(fmt.Sprintf("The following keys are missing: %s", list))
}
