/**
 * Date: 23.08.12
 * Time: 23:47
 *
 * @author Vladimir Matveev
 */
package common

import (
    msg "bridge/common/msg"
)

type LocalPlugin interface {
    SupportsCommand(command string) bool
    HandleMessage(input msg.Message) msg.Message
}

type RemotePlugin interface {
    SupportsCommand(command string) bool
    HandleMessage(input msg.Message) msg.Message
}
