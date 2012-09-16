/**
 * Date: 16.09.12
 * Time: 12:55
 *
 * @author Vladimir Matveev
 */
package repo

import (
    . "bridge/common"
)

type MediatorMaker func () Mediator

var mediatorsRepo = map[string]MediatorMaker {
}

func GetMediator(name string) Mediator {
    mmaker := mediatorsRepo[name]
    if mmaker == nil {
        return nil
    }
    return mmaker()
}
