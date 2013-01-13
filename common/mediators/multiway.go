/**
 * Date: 16.09.12
 * Time: 16:43
 *
 * @author Vladimir Matveev
 */
package mediators

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
    "sync"
    "runtime"
)

type multiway struct {
    endpoints map[string]bool
    ports map[string][]Port
    terminated bool
    lock sync.Mutex
}

func NewMultiway() Mediator {
    return &multiway{
        endpoints: make(map[string]bool),
        ports: make(map[string][]Port),
        terminated: false,
    }
}

func (m *multiway) Name() string {
    return "multiway"
}

func (m *multiway) Config(mconf *conf.MediatorConf) error {
    if string(mconf.Mediator) != m.Name() {
        return newErrorf("Invalid mediator name: %s", mconf.Mediator)
    }

    if len(mconf.EndpointNames) == 0 {
        return newErrorf("Endpoint names are not specified")
    } else {
        for _, endpoint := range mconf.EndpointNames {
            m.endpoints[endpoint] = true
        }
    }

    go m.notifier()

    return nil
}

func (m *multiway) HasEndpoint(endpoint string) bool {
    return m.endpoints[endpoint]
}

func (m *multiway) Connect(endpoint string, port Port) error {
    // Take a hold on the mutex
    m.lock.Lock()
    defer m.lock.Unlock()

    if !m.HasEndpoint(endpoint) {
        return newErrorf("Invalid endpoint name: %s", endpoint)
    }

    if ports, ok := m.ports[endpoint]; ok {
        m.ports[endpoint] = append(ports, port)
    } else {
        ports = make([]Port, 1, 4)
        ports[0] = port
        m.ports[endpoint] = ports
    }

    return nil
}

func (m *multiway) Term() {
    m.terminated = true
}

func (m *multiway) notifier() {
    // Prepare a slice where we will hold ports
    ports := make([]Port, 0, len(m.ports))

    // Main loop
    for {
        // First, check whether we have been interrupted
        if m.terminated {
            return
        }

        // We need exclusive access to ports map, so we can construct consistent slice of ports

        // Lock the mutex
        m.lock.Lock()

        // Clear the ports slice
        ports = ports[:0]

        // Fill the ports slice
        // TODO: add 'changed' flag and refill the slice only if it is set
        for _, endpointPorts := range m.ports {
            for _, port := range endpointPorts {
                ports = append(ports, port)
            }
        }

        // Unlock the mutex back
        m.lock.Unlock()

        // Walk through all ports and fan out all incoming messages to other ports
        for _, port := range ports {
            select {
            case msg, ok := <-port.Second().Source():
                // If the receive was successful, fan the received message out to all other ports
                if ok {
                    for _, other := range ports {
                        if other != port {
                            other.Second().Sink() <- msg
                        }
                    }
                }
            default:
                // Pass to the next port
            }
        }

        // Give time to other routines
        runtime.Gosched()
    }
}
