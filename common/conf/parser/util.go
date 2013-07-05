/**
 * Date: 14.01.13
 * Time: 23:37
 *
 * @author Vladimir Matveev
 */
package parser

import (
    "code.google.com/p/gelo"
    "fmt"
    "strings"
)

func runtimeError(vm *gelo.VM, fmt string, args ...interface{}) {
    gelo.RuntimeError(vm, fmt.Sprintf(fmt, args...))
}

func getOrMakeDict(vm *gelo.VM, name string) *gelo.Dict {
    k := gelo.StrToSym(name)
    dw, ok := vm.Ns.Get(0, k)
    var d *gelo.Dict
    if ok {
        d, ok = dw.(*gelo.Dict)
        if !ok {
            d = gelo.NewDict()
            vm.Ns.Set(0, k, d)
        }
    } else {
        d = gelo.NewDict()
        vm.Ns.Set(0, k, d)
    }
    return d
}

func checkInSection(vm *gelo.VM, command string, sections ...string) string {
    if len(sections) == 0 {
        checkNotInSection(vm, command)
        return ""
    }

    tag, name := currentSection(vm)
    for section := range sections {
        if tag == section {
            return name
        }
    }

    invalidSectionError(vm, command, strings.Join(sections, " or "), tag)
    return ""
}

func checkNotInSection(vm *gelo.VM, command string) {
    k := gelo.StrToSym("section")
    if s, ok := vm.Ns.Get(0, k); ok {
        gelo.RuntimeError(
            vm,
            fmt.Sprintf(
                "%s command is used inside a section %s while it should be toplevel",
                command, s.Ser().String(),
            ),
        )
    }
}

func invalidSectionError(vm *gelo.VM, command string, required string, got string) {
    runtimeError(vm, "%s command is used in invalid section %s instead of %s", command, got, required)
}

func enterSection(vm *gelo.VM, tag string, name string) {
    vm.Ns.Fork(nil)

    k := gelo.StrToSym("section-tag")
    v := gelo.StrToSym(tag)
    if !vm.Ns.Set(0, k, v) {
        runtimeError(vm, "Unable to set section tag: section %s %s", tag, name)
    }

    k = gelo.StrToSym("section-name")
    v = gelo.StrToSym(name)
    if !vm.Ns.Set(0, k, v) {
        runtimeError(vm, "Unable to set section name: section %s %s", tag, name)
    }
}

func enterNamelessSection(vm *gelo.VM, tag string) {
    enterSection(vm, tag, "")
}

func leaveSection(vm *gelo.VM, tag string, name string) {
    // In fact deleting keys is completely unneccessary; it is added only for additional check that we really are
    // in a section
    k := gelo.StrToSym("section-tag")
    if _, ok := vm.Ns.Del(k); !ok {
        runtimeError(vm, "Unable to delete section tag (not in the section?): section %s %s", tag, name)
    }
    k := gelo.StrToSym("section-name")
    if _, ok := vm.Ns.Del(k); !ok {
        runtimeError(vm, "Unable to delete section name (not in the section?): section %s %s", tag, name)
    }

    k = gelo.StrToSym
    if !vm.Ns.Unfork() {
        runtimeError(vm, "Unable to exit the namespace (not in the section?): section %s %s", tag, name)
    }
}

func leaveNamelessSection(vm *gelo.VM, tag string) {
    leaveSection(vm, tag, "")
}

func currentSection(vm *gelo.VM) (string, string) {
    k := gelo.StrToSym("section-tag")
    tag, ok := vm.Ns.Get(0, k)
    if !ok {
        runtimeError(vm, "Unable to get current section tag (not in the section?)")
    }

    k = gelo.StrToSym("section-name")
    name, ok := vm.Ns.Get(0, k)
    if !ok {
        runtimeError(vm, "Unable to get current section name (not in the section?)")
    }

    return tag.Ser().String(), name.Ser().String()
}

func insideSection(vm *gelo.VM, tag string, name string, action func()) {
    enterSection(vm, tag, name)
    action()
    leaveSection(vm, tag, name)
}

func insideNamelessSection(vm *gelo.VM, tag string, action func()) {
    insideSection(vm, tag, "", action)
}
