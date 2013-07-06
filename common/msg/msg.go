/**
 * Date: 24.08.12
 * Time: 16:29
 *
 * @author Vladimir Matveev
 */
package msg

import (
    "code.google.com/p/goprotobuf/proto"
    "encoding/binary"
    "fmt"
    "io"
)

type SerializeError struct {
    msg   string
    cause error
}

func (err *SerializeError) Error() string {
    if err.cause == nil {
        return err.msg
    }
    return err.msg + ": " + err.cause.Error()
}

func (err *SerializeError) Cause() error {
    return err.cause
}

func NewError(cause error, msg string, args ...interface{}) error {
    return &SerializeError{fmt.Sprintf(msg, args...), cause}
}

func (msg *Message) Header(name string) *string {
    for _, h := range msg.Headers {
        if h.GetName() == name {
            return &h.GetValue()
        }
    }
    return nil
}

func (msg *Message) SetHeader(name string, value string) {
    for _, h := range msg.Headers {
        if h.GetName() == name {
            h.Value = proto.String(value)
            return
        }
    }
    hdr := &Header{Name: proto.String(name), Value: proto.String(value)}
    msg.Headers = append(msg.Headers, hdr)
}

func (msg *Message) BodyPart(name string) []byte {
    for _, b := range msg.BodyParts {
        if b.GetName() == name {
            return b.GetData()
        }
    }
    return nil
}

func (msg *Message) SetBodyPart(name string, data []byte) {
    for _, b := range msg.BodyParts {
        if b.GetName() == name {
            b.Data = data
            return
        }
    }
    bp := &BodyPart{Name: proto.String(name), Data: data}
    msg.BodyParts = append(msg.BodyParts, bp)
}

func CreateEmpty() *Message {
    return CreateWithName("")
}

func CreateWithName(name string) *Message {
    msg := new(Message)
    msg.Name = proto.String(name)
    return msg
}

func Serialize(msg *Message, w io.Writer) error {
    data, err := proto.Marshal(msg)
    if err != nil {
        return NewError(err, "Error serializing message")
    }

    size := make([]byte, 4)
    binary.LittleEndian.PutUint32(size, uint32(len(data)))
    n, err := w.Write(size)
    if err != nil {
        return NewError(err, "Error writing message size to output, wrote %d bytes", n)
    }

    n, err = w.Write(data)
    if err != nil {
        return NewError(err, "Error writing message content to output, wrote %d bytes", n)
    }

    return nil
}

func Deserialize(r io.Reader) (*Message, error) {
    sizeBuf := make([]byte, 4)
    if n, err := io.ReadAtLeast(r, sizeBuf, 4); err != nil {
        return nil, NewError(err, "Error reading message size from input, read %d bytes", n)
    }
    size := binary.LittleEndian.Uint32(sizeBuf)

    data := make([]byte, size)
    if n, err := io.ReadAtLeast(r, data, int(size)); err != nil {
        return nil, NewError(err, "Error reading message content from input, read %d bytes", n)
    }

    msg := new(Message)
    if err := proto.Unmarshal(data, msg); err != nil {
        return nil, NewError(err, "Error deserializing message")
    }

    return msg, nil
}
