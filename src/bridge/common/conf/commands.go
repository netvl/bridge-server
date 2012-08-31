package conf

import (
    "code.google.com/p/gelo"
)

var confCommands = map[string]interface{} {
    "conf": conf,
    "set": set,
}

func Conf(vm *gelo.VM, args *gelo.List, ac uint) gelo.Word {
    if ac != 2 {
        gelo.ArgumentError(vm, "conf", "name body", args)
    }
    name := args.Value
    body := vm.API.QuoteOrElse(args.Next.Value)


}
