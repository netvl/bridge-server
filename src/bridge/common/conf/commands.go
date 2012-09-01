package conf

import (
    "code.google.com/p/gelo"
    "log"
)

var confCommands = map[string]interface{}{
    "conf": conf,
    "set":  set,
}

func conf(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 2 {
        gelo.ArgumentError(vm, "conf", "name body", args)
    }
    name := vm.API.SymbolOrElse(args.Value)
    body := vm.API.QuoteOrElse(args.Next.Value)
    log.Println("In conf:", name)

    d := getOrMakeDict(vm)

    // The value returned by inner command will be dict itself
    vm.Ns.Fork(nil)
    value := vm.API.InvokeCmdOrElse(body, nil)
    vm.Ns.Unfork()
    d.Set(name, value)

    log.Println("Returning from conf", name)
    return d
}

func set(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac == 0 {
        gelo.ArgumentError(vm, "set", "name [args*]", args)
    }
    name := vm.API.SymbolOrElse(args.Value)
    log.Println("In set:", name)

    // Check that all arguments are symbols
    for tmp := args; tmp != nil; tmp = tmp.Next {
        _, oksym := tmp.Value.(gelo.Symbol)
        _, oknum := tmp.Value.(*gelo.Number)
        _, okbool := tmp.Value.(gelo.Bool)
        if !oksym && !oknum && !okbool {
            gelo.RuntimeError(vm, "Arguments should be symbols or numbers")
        }
    }

    var values *gelo.List
    if ac == 1 {
        values = gelo.NewListFromGo([]interface{}{gelo.True})
    } else {
        values = args.Next
    }

    d := getOrMakeDict(vm)
    d.Set(name, values)

    log.Println("Returning from set", name)
    return d
}

func getOrMakeDict(vm *gelo.VM) *gelo.Dict {
    k := gelo.StrToSym("data")
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
