package conf

import (
    "code.google.com/p/gelo"
    "log"
    "os"
)

func LoadConfig(file string) *Conf {
    // gelo.SetTracer(extensions.Stderr)
    // gelo.TraceOn(gelo.All_traces)

    port := gelo.NewChan()
    vm := gelo.NewVM(port)
    defer vm.Destroy()

    vm.RegisterBundle(confCommands)

    f, err := os.Open(file)
    defer f.Close()
    if err != nil {
        log.Printf("Error opening config file: %v", err)
        return nil
    }
    result, err := vm.Run(f, nil)
    if err != nil {
        log.Printf("VM error: %v", err)
    } else {
        log.Println(result)
    }

    return nil
}
