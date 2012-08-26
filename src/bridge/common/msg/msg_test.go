/**
 * Date: 25.08.12
 * Time: 20:05
 *
 * @author Vladimir Matveev
 */
package msg

import (
    . "launchpad.net/gocheck"
    "testing"
    "bytes"
    "encoding/hex"
)

const (
    EXAMPLE1_HEXDUMP = "4d5347530000000c746573745f6d657373616765020468647231060076616c756531046864723206" +
        "0076616c7565320203627031040000000102030403627032160000005465737420626f6479207061727420636f6e74656e74"
)

func createMsg_Example1() *Message {
    msg := CreateWithName("test_message")
    msg.SetHeader("hdr1", "value1")
    msg.SetHeader("hdr2", "value2")
    msg.SetBodyPart("bp1", BodyPartFromSlice([]byte{1, 2, 3, 4}))
    msg.SetBodyPart("bp2", BodyPartFromSlice([]byte("Test body part content")))
    return msg
}

func Test(t *testing.T) {
    TestingT(t)
}

type MsgSuite struct {}
var _ = Suite(&MsgSuite{})

func (s *MsgSuite) TestMessageSerialization(t *C) {
    msg := createMsg_Example1()

    var buf bytes.Buffer
    if err := Serialize(msg, &buf); err != nil {
        t.Fatal("Serialization error: " + err.Error())
    }

    hexdump := hex.EncodeToString(buf.Bytes())
    if hexdump != EXAMPLE1_HEXDUMP {
        t.Fatal("Serialized messages do not match!")
    }
}

func (s *MsgSuite) TestMessageDeserialization(t *C) {
    data, _ := hex.DecodeString(EXAMPLE1_HEXDUMP)
    src := bytes.NewReader(data)

    msg, err := DeserializeMessageName(src)
    t.Assert(err, IsNil)

    err = DeserializeMessageHeaders(src, msg)
    t.Assert(err, IsNil)

    err = DeserializeMessageBodyParts(src, msg, EmptyHook)
    t.Assert(err, IsNil)

    tpl := createMsg_Example1()

    t.Assert(msg.GetName(), Equals, tpl.GetName())
    t.Assert(msg.headers, DeepEquals, tpl.headers)
    t.Assert(msg.bodyParts, DeepEquals, tpl.bodyParts)  // TODO: fix it
}
