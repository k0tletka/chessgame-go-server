package util

import (
    "regexp"
    "reflect"

    pcre "github.com/gijsbers/go-pcre"
)

var (
    // Pcre precompiled patterns
    passwordPattern = pcre.MustCompile(`(?=.*?[0-9])(?=.*?[A-Z])[a-zA-Z0-9@#!$%&_]{8,}`, pcre.ANCHORED)

    // Map of validation functions
    validatorFuncs = map[string]func(interface{}) bool{
        "Login":
            func(value interface{}) bool {
                login, ok := value.(string)

                if !ok {
                    return false
                }

                matched, err := regexp.MatchString(`^[a-z0-9]{6,100}$`, login)
                return err == nil && matched
            },
        "Password":
            func(value interface{}) bool {
                password, ok := value.(string)

                if !ok {
                    return false
                }

                return passwordPattern.MatcherString(password, pcre.ANCHORED).Matches()
            },
        "Email":
            func(value interface{}) bool {
                email, ok := value.(string)

                if !ok {
                    return false
                }

                matched, err := regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`, email)
                return err == nil && matched
            },
        "NotDefaultValue":
            func(value interface{}) bool {
                return reflect.ValueOf(value).IsZero()
            },
    }
)

// Type of validating value
type VValue struct {
    Type string
    Value interface{}
}

// Validation function
func ValidateValues(values... *VValue) bool {
    for _, v := range values {
        // Get validation func by passed type
        if validateFunc, ok := validatorFuncs[v.Type]; ok {
            return validateFunc(v.Value)
        }

        return false
    }

    return true
}
