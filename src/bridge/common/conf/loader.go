package conf

import (
    "code.google.com/p/gelo"
    "code.google.com/p/gelo/commands"
)

func LoadConfig(file string) {
    port := gelo.NewChan()
    vm := gelo.NewVM(port)
    vm.RegisterBundle(confCommands)
}
