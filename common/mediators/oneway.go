/**
 * Date: 17.09.12
 * Time: 10:12
 *
 * @author Vladimir Matveev
 */
package mediators

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/conf"
)

type oneway struct {
    endpoint string
    subscribers []Subscriber
}

func NewOneway() Mediator {
    return &oneway{
        subscribers: make([]Subscriber, 0, 2),
    }
}

func (o *oneway) Name() string {
    return "oneway"
}

func (o *oneway) Config(mconf *conf.MediatorConf) error {
    if string(mconf.Mediator) != o.Name() {
        return newErrorf("Invalid mediator name: %s", mconf.Mediator)
    }

    if len(mconf.EndpointNames) == 0 {
        return newErrorf("Endpoint name is not specified")
    } else {
        o.endpoint = mconf.EndpointNames[0]
    }

    return nil
}

func (o *oneway) Submit(endpoint string, msg interface{}) error {
    if endpoint != o.endpoint {
        return newErrorf("Invalid endpoint %s requested", endpoint)
    }

    for _, s := range o.subscribers {
        s(msg)
    }

    return nil
}

func (o *oneway) Subscribe(endpoint string, s Subscriber) error {
    if endpoint != o.endpoint {
        return newErrorf("Invalid endpoint %s requested", endpoint)
    }

    if o.subscribers == nil {
        o.subscribers = make([]Subscriber, 0, 1)
    }

    o.subscribers = append(o.subscribers, s)

    return nil
}

func (o *oneway) Term() {
    // Nothing
}
