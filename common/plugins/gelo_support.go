/**
 * Date: 04.07.13
 * Time: 21:38
 *
 * @author Vladimir Matveev
 */
package plugins

import (
    "code.google.com/p/gelo"
    "fmt"
)

type SenderType int
type SenderName string

const (
    SENDER_PLUGIN SenderType = iota
    SENDER_COMMUNICATOR
)

func (t SenderType) String() string {
    switch t {
    case SENDER_PLUGIN:
        return "Plugin"
    case SENDER_COMMUNICATOR:
        return "Communicator"
    }
    return nil
}

func (t SenderName) String() string {
    return string(t)
}

type SenderIdentity struct {
    senderType SenderType
    senderName SenderName
}

func (id SenderIdentity) SenderType() SenderType {
    return id.senderType
}

func (id SenderIdentity) SenderName() SenderName {
    return id.senderName
}

func (id SenderIdentity) Type() gelo.Symbol {
    return gelo.StrToSym("SenderIdentity")
}

func (id SenderIdentity) Ser() gelo.Symbol {
    return gelo.StrToSym(fmt.Sprintf("%s %s", id.senderType, id.senderName))
}

func (id SenderIdentity) Copy() gelo.Word {
    return id
}

func (id SenderIdentity) DeepCopy() gelo.Word {
    return id
}

func (id SenderIdentity) Equals(word gelo.Word) bool {
    if id2, ok := word.(SenderIdentity); ok {
        return int(id.senderType) == id(id2.senderType) && string(id.senderName) == string(id2.senderName)
    }
    return false
}
