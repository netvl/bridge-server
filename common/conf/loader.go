package conf

import (
    "code.google.com/p/gelo"
    "io"
    "log"
    "os"
)

func LoadConfigFromReader(src io.Reader) (*Conf, error) {
    port := gelo.NewChan()
    vm := gelo.NewVM(port)
    defer vm.Destroy()

    vm.RegisterBundle(confCommands)

    result, err := vm.Run(src, nil)
    if err != nil {
        return nil, err
    }

    rd, _ := result.(*gelo.Dict)

    return buildConfig(rd)
}

// Returns a configuration loaded from a file with given name.
func LoadConfigFromFile(file string) (*Conf, error) {
    f, err := os.Open(file)
    if err != nil {
        log.Printf("Error opening config file: %v", err)
        return nil, err
    }
    defer f.Close()

    return LoadConfigFromReader(f)

}
