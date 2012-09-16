/**
 * Date: 16.09.12
 * Time: 16:43
 *
 * @author Vladimir Matveev
 */
package mediators

import (
    . "bridge/common"
    "bridge/common/conf"
    "strconv"
)

type envelope struct {
    endpoint string
    datum interface{}
}

type multiway struct {
    container chan *envelope
    endpoints map[string]bool
    subscribers map[string][]Subscriber
}

func NewMultiway() Mediator {
    return &multiway{
        endpoints: make(map[string]bool),
        subscribers: make(map[string][]Subscriber),
    }
}

func (m *multiway) Name() string {
    return "multiway"
}

func (m *multiway) Config(mconf *conf.MediatorConf) error {
    if string(mconf.Mediator) != m.Name() {
        return newError("Invalid mediator name: %s", mconf.Mediator)
    }

    var capacity string
    if capacitySlice, ok := mconf.Options["capacity"]; ok && len(capacitySlice) > 0 {
        capacity = capacitySlice[0]
    } else {
        capacity = "infinite"
    }

    if capacity == "infinite" {
        m.container = make(chan *envelope)
    } else if n, err := strconv.Atoi(capacity); err != nil {
        return newError("Capacity value '%s': %s", capacity, err)
    } else {
        m.container = make(chan *envelope, n)
    }

    if len(mconf.EndpointNames) == 0 {
        return newError("Endpoint name is not specified")
    } else {
        for _, endpoint := range mconf.EndpointNames {
            m.endpoints[endpoint] = true
        }
    }

    go m.notifier()

    return nil
}

func (m *multiway) Submit(endpoint string, msg interface{}) error {
    if !m.endpoints[endpoint] {
        return newError("Invalid endpoint %s requested", endpoint)
    }
    m.container <- &envelope{endpoint, msg}

    return nil
}

func (m *multiway) Subscribe(endpoint string, s Subscriber) error {
    if !m.endpoints[endpoint] {
        return newError("Invalid endpoint %s requested", endpoint)
    }

    if subscribers, ok := m.subscribers[endpoint]; ok {
        m.subscribers[endpoint] = append(subscribers, s)
    } else {
        subscribers = make([]Subscriber, 1, 4)
        subscribers[0] = s
        m.subscribers[endpoint] = subscribers
    }

    return nil
}

func (m *multiway) Term() {
    close(m.container)
}

func (m *multiway) notifier() {
    for {
        msg, ok := <-m.container
        if !ok {
            break
        }

        for ep, ss := range m.subscribers {
            if ep == msg.endpoint {
                continue
            }
            for _, s := range ss {
                s(msg.datum)
            }
        }
    }
}
