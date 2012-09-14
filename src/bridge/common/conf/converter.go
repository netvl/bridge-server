/**
 * Date: 03.09.12
 * Time: 11:30
 *
 * @author Vladimir Matveev
 */
package conf

import "code.google.com/p/gelo"

// configElement represents single config element, that is, a map or a slice of string values.
type configElement interface {
    AsDict() map[string]configElement
    AsEntry() []string
}

// configDict is a shortcut for the map of configElements.
type configDict map[string]configElement

// An implementation of a method from configElement interface, returns underlying map.
func (d configDict) AsDict() map[string]configElement {
    return map[string]configElement(d)
}

// An implementation of a method from configElement interface, panic with *ConfigError.
func (d configDict) AsEntry() []string {
    panic(newConfigError("", "Attempt to use dict as config entry: %v", d))
    return nil
}

// configEntry is a shortcut for a slice of string values, represents configuration value.
type configEntry []string

// An implementation of a method from configElement interface, panic with *ConfigError.
func (e configEntry) AsDict() map[string]configElement {
    panic(newConfigError("", "Attempt to use config entry as dict: %v", e))
    return nil
}

// An implementation of a method from configElement interface, returns underlying slice.
func (e configEntry) AsEntry() []string {
    return []string(e)
}

// convertDict converts Gelo dictionary to a map from strings to configElements,
// recursively converting elements of the dictionary.
func convertDict(d *gelo.Dict) map[string]configElement {
    rm := make(map[string]configElement)

    m := d.Map()
    for k, w := range m {
        switch cw := w.(type) {
        case *gelo.Dict:
            cm := convertDict(cw)
            rm[k] = configDict(cm)
        case *gelo.List:
            ce := convertList(cw)
            rm[k] = configEntry(ce)
        default:
            panic(newConfigError("", "Illegal object %v encountered at key %v inside %v dict", w, k, d))
        }
    }

    return rm
}

// convertList converts Gelo list to a slice of string values, serializing them first.
func convertList(l *gelo.List) []string {
    rl := make([]string, 0, 1)
    for ; l != nil; l = l.Next {
        rl = append(rl, l.Value.Ser().String())
    }
    return rl
}

// convertDictSafe wraps convertDict, recovering from any panicks it may have started,
// returning an error as plain value.
func convertDictSafe(d *gelo.Dict) (cm map[string]configElement, err error) {
    defer func() {
        v := recover()
        if v != nil {
            if e, ok := v.(*ConfigError); ok {
                err = e
            } else {
                panic(v)
            }
        }
    }()

    return convertDict(d), nil
}

