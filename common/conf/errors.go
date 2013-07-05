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

// ConfigError is a configuration error, that is, some error occured at a specific
// location inside the configuration file.
type ConfigError struct {
    Msg      string
    Location string
}

// Returns a new configuration error with given location and formatted message.
func newConfigError(location string, msg string, args ...interface{}) *ConfigError {
    return &ConfigError{fmt.Sprintf(msg, args...), location}
}

// Error is an implementation of the error interface, returns the message and
// the location at which the error occured.
func (err *ConfigError) Error() string {
    return err.Msg + " @ " + err.Location
}

// ConfigErrors encapsulates a slice of ConfigError values. Used to incrementally build
// a list of errors when building configuration object.
type ConfigErrors struct {
    errors []*ConfigError
}

// Returns new empty ConfigErrors pointer.
func makeConfigErrors() *ConfigErrors {
    ce := make([]*ConfigError, 0, 2)
    return &ConfigErrors{ce}
}

// Returns true if there are no errors are present in the container, false otherwise.
func (errs *ConfigErrors) noErrors() bool {
    return len(errs.errors) == 0
}

// Add adds an error with given location and formatted message to the list of errors,
// returning this list.
func (errs *ConfigErrors) add(location string, msg string, args ...interface{}) *ConfigErrors {
    errs.errors = append(errs.errors, newConfigError(location, msg, args...))
    return errs
}

// AddError adds an error value to the list of errors, returning this list.
func (errs *ConfigErrors) addError(err *ConfigError) *ConfigErrors {
    errs.errors = append(errs.errors, err)
    return errs
}

// Merge adds all entries from the provided list to the receiving list,
// returning the latter one.
func (errs *ConfigErrors) merge(other *ConfigErrors) *ConfigErrors {
    errs.errors = append(errs.errors, other.errors...)
    return errs
}

// PrependLocation adds same prefix to all of configuration errors locations.
// If a location in some of the errors is empty or spaces-only string, it will be ignored.
func (errs *ConfigErrors) prependLocation(location string) {
    for _, ce := range errs.errors {
        if strings.TrimSpace(ce.Location) == "" {
            ce.Location = location
        } else {
            ce.Location = location + ", " + ce.Location
        }
    }
}

// String returns newline-joined list of strings, each is formed by calling Error method
// on the next error in the list.
func (errs *ConfigErrors) String() string {
    msgs := make([]string, 0, 2)
    for _, ce := range errs.errors {
        msgs = append(msgs, ce.Error())
    }
    return strings.Join(msgs, "\n")
}

// Error is the same as String in this case.
func (errs *ConfigErrors) Error() string {
    return errs.String()
}

// Returns an underlying slice of configuration errors.
func (errs *ConfigErrors) Errors() []*ConfigError {
    return errs.Errors()
}
