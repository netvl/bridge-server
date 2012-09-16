/**
 * Date: 16.09.12
 * Time: 17:07
 *
 * @author Vladimir Matveev
 */
package mediators

import "fmt"

type MediatorError struct {
    Msg string
}

func (e *MediatorError) Error() string {
    return e.Msg
}

func newError(format string, args ...interface{}) error {
    return &MediatorError{fmt.Sprintf(format, args...)}
}

