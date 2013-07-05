/**
 * Date: 14.01.13
 * Time: 2:36
 *
 * @author Vladimir Matveev
 */
package util

func StringPattern(args ...interface{}) []*string {
    result := make([]*string, len(args))
    for i, arg := range args {
        if arg == nil {
            result[i] = nil
        } else if p, ok := arg.(*string); ok {
            result[i] = p
        } else {
            result[i] = new(string)
            *result[i] = arg.(string)
        }
    }
    return result
}

func MatchStringSlicePattern(source []string, args ...interface{}) bool {
    pattern := StringPattern(args...)

    pl, sl := len(pattern), len(source)

    if pattern == nil {
        // Nil pattern matches anything
        return true
    } else if len(pattern) == 0 {
        // Empty pattern matches nothing
        return false
    } else if pl-sl > 1 || (pl-sl == 1 && pattern[pl-1] != nil) {
        // If source is shorter than pattern by more than 1 element, or if it is shorter by 1 element and pattern's
        // last element is not nil, then there is no match
        return false
    } else if sl > pl && pattern[pl-1] != nil {
        // If source is greater than pattern and pattern's last element is not nil, then there is no match
        return false
    }

    min := len(pattern)
    if s := len(source); s < min {
        min = s
    }

    for i := 0; i < min; i++ {
        p, s := pattern[i], source[i]

        // Nil pattern means match anything
        if p == nil {
            continue
        }

        // Otherwise check whether pattern and source elements are equal
        if *p != s {
            return false
        }
    }

    if sl == pl || sl >= pl-1 && pattern[pl-1] == nil {
        return true
    }

    return false
}
