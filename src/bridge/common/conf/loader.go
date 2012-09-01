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
    if err != nil {
        log.Printf("Error opening config file: %v", err)
        return nil
    }
    defer f.Close()

    result, err := vm.Run(f, nil)
    if err != nil {
        log.Printf("VM error: %v", err)
        return nil
    }

    rd, _ := result.(*gelo.Dict)

    return buildConfig(rd)
}
