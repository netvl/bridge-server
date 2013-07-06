/**
 * Date: 06.07.13
 * Time: 17:54
 *
 * @author Vladimir Matveev
 */
package util

// Performs bitwise AND operation on bytes in a and b slices and returns the result as a new slice.
// If slices are of unequal length, nil is returned.
func AndSlices(a, b []byte) (result []byte) {
    if len(a) != len(b) {
        return
    }

    result = make([]byte, len(a))

    for i := 0; i < len(a); i++ {
        result[i] = a[i] & b[i]
    }

    return
}

// Performs bitwise OR operation on bytes in a and b slices and returns the result as a new slice.
// If slices are of unequal length, nil is returned.
func OrSlices(a, b []byte) (result []byte) {
    if len(a) != len(b) {
        return
    }

    result = make([]byte, len(a))

    for i := 0; i < len(a); i++ {
        result[i] = a[i] | b[i]
    }

    return
}

// Performs bitwise NOT operation on bytes in a slice and returns the result as a new slice.
func NotSlice(a []byte) (result []byte) {
    result = make([]byte, len(a))

    for i := 0; i < len(a); i++ {
        result[i] = ^a[i]
    }

    return
}
