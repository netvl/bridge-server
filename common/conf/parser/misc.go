/**
 * Date: 14.01.13
 * Time: 23:39
 *
 * @author Vladimir Matveev
 */
package parser

import "code.google.com/p/gelo"

func booleanOption(vm *gelo.VM, args *gelo.List, ac uint, dst *bool, command string, sections ...string) {
    if ac > 1 {
        gelo.ArgumentError(vm, command, "[true | false]", args)
    }

    // Check whether we're in the specified sections
    checkInSection(vm, command, sections...)

    if ac > 0 {
        enabled, ok := args.Value.(gelo.Bool)
        if !ok {
            gelo.RuntimeError(vm, command, "argument should be equal to true or false or be absent")
        }
        *dst = enabled.True()
    } else {
        *dst = true
    }
}

func (p *ConfigParser) set(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac == 0 {
        gelo.ArgumentError(vm, "set", "<name> [args*]", args)
    }
    name := vm.API.SymbolOrElse(args.Value)

    // Check that we're in correct section
    checkInSection(vm, "set", "port", "plugin")

    // Check that all arguments are compatible values
    for tmp := args; tmp != nil; tmp = tmp.Next {
        _, oksym := tmp.Value.(gelo.Symbol)
        _, oknum := tmp.Value.(*gelo.Number)
        _, okbool := tmp.Value.(gelo.Bool)
        if !(oksym || oknum || okbool) {
            runtimeError(vm, "Arguments of set should be symbols or numbers or booleans")
        }
    }

    var values *gelo.List
    if ac == 1 {
        values = gelo.NewListFromGo([]interface{}{gelo.True})
    } else {
        values = args.Next
    }

    d := getOrMakeDict(vm, "data")
    d.Set(name, values)

    return nil
}
