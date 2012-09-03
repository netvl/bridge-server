/**
 * Date: 03.09.12
 * Time: 9:05
 *
 * @author Vladimir Matveev
 */
package conf

import (
    "fmt"
    "strings"
)

type ConfigError struct {
    Msg string
    Location string
}

func newConfigError(location string, msg string, args ...interface{}) *ConfigError {
    return &ConfigError{fmt.Sprintf(msg, args...), location}
}

func (err *ConfigError) Error() string {
    return err.Msg
}

type ConfigErrors struct {
    errors []*ConfigError
}

func makeConfigErrors() *ConfigErrors {
    ce := make([]*ConfigError, 0, 2)
    return &ConfigErrors{ce}
}

func (errs *ConfigErrors) Add(location string, msg string, args ...interface{}) *ConfigErrors {
    errs.errors = append(errs.errors, newConfigError(location, msg, args...))
    return errs
}

func (errs *ConfigErrors) AddError(err *ConfigError) *ConfigErrors {
    errs.errors = append(errs.errors, err)
    return errs
}

func (errs *ConfigErrors) Merge(other *ConfigErrors) *ConfigErrors {
    errs.errors = append(errs.errors, other.errors...)
    return errs
}

func (errs *ConfigErrors) PrependLocation(location string) {
    for _, ce := range errs.errors {
        if strings.TrimSpace(ce.Location) == "" {
            ce.Location = location
        } else {
            ce.Location = location + ", " + ce.Location
        }
    }
}

func (errs *ConfigErrors) String() string {
    msgs := make([]string, 0, 2)
    for _, ce := range errs.errors {
        msgs = append(msgs, fmt.Sprintf("%s @ %s", ce.Msg, ce.Location))
    }
    return strings.Join(msgs, "\n")
}


