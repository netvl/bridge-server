/**
 * Date: 16.09.12
 * Time: 12:55
 *
 * @author Vladimir Matveev
 */
package repo

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "github.com/dpx-infinity/bridge-server/common/mediators"
)

type MediatorMaker func () Mediator

var mediatorsRepo = map[string]MediatorMaker {
    "multiway": func() Mediator { return mediators.NewMultiway() },
    "oneway": func() Mediator { return mediators.NewOneway() },
}

func GetMediator(name string) Mediator {
    mmaker := mediatorsRepo[name]
    if mmaker == nil {
        return nil
    }
    return mmaker()
}
