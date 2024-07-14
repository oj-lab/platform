package casbin_agent

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// KeyMatchGin determines whether key1 matches the pattern of key2 (similar to KeyMatch2 in Casbin).
// Providing a Gin-like path matching pattern.
func KeyMatchGin(key1 string, key2 string) bool {
	// /user/:name can match /user/john but not /user/ or /user/john/
	re := regexp.MustCompile(`:[^/]+`)
	key2 = re.ReplaceAllString(key2, `([^/]+)`)
	// /user/:name/*action can match /user/john/ and /user/john/send and /user/john/send/email
	re = regexp.MustCompile(`/\*[^/]+`)
	key2 = re.ReplaceAllString(key2, `(/.*|.{0})`)
	// Turn every / into \/ and add ^ and $ to match the whole string
	key2 = ("^" + strings.ReplaceAll(key2, "/", "\\/") + "$")

	return regexp.MustCompile(key2).MatchString(key1)
}

func keyMatchGinFunc(args ...interface{}) (interface{}, error) {
	if err := validateVariadicArgs(2, args...); err != nil {
		return false, fmt.Errorf("%s: %s", "keyMatch2", err)
	}

	name1 := args[0].(string)
	name2 := args[1].(string)

	return bool(KeyMatchGin(name1, name2)), nil
}

// Copied from casbin/util/builtin_operators.go.
// Validate the variadic parameter size and type as string
func validateVariadicArgs(expectedLen int, args ...interface{}) error {
	if len(args) != expectedLen {
		return fmt.Errorf("expected %d arguments, but got %d", expectedLen, len(args))
	}

	for _, p := range args {
		_, ok := p.(string)
		if !ok {
			return errors.New("argument must be a string")
		}
	}

	return nil
}
