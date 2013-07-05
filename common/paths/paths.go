/**
 * Date: 05.07.13
 * Time: 20:49
 *
 * @author Vladimir Matveev
 */
package paths

import "os"

const (
    def_lib = "/usr/lib/bridge"
    def_plugins = "/plugins"
)

const (
    ENV_LIB = "BRIDGE_LIB"
    ENV_PLUGINS = "BRIDGE_PLUGINS"
)

var (
    lib string
    plugins string
)

func init() {
    if envLib := os.Getenv(ENV_LIB); envLib == "" {
        lib = def_lib
    } else {
        lib = envLib
    }

    if envPlugins := os.Getenv(ENV_PLUGINS); envPlugins == "" {
        plugins = lib + def_plugins
    } else {
        plugins = envPlugins
    }
}

// Returns computed location of bridge lib directory.
func Lib() string {
    return lib
}

// Returns computed location of bridge plugins directory.
func Plugins() string {
    return plugins
}
