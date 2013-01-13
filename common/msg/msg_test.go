/**
 * Date: 25.08.12
 * Time: 20:05
 *
 * @author Vladimir Matveev
 */
package msg

import (
    "bytes"
    . "launchpad.net/gocheck"
    "testing"
)

func createMsg_Example1() *Message {
    msg := CreateWithName("test_message")
    msg.SetHeader("hdr1", "value1")
    msg.SetHeader("hdr2", "value2")
    msg.SetBodyPart("bp1", []byte{1, 2, 3, 4})
    msg.SetBodyPart("bp2", []byte("Test body part content"))
    return msg
}

func Test(t *testing.T) {
    TestingT(t)
}

type MsgSuite struct{}

var _ = Suite(&MsgSuite{})

func (s *MsgSuite) TestMessageSerializationDeserialization(t *C) {
    msg := createMsg_Example1()

    var buf bytes.Buffer
    if err := Serialize(msg, &buf); err != nil {
        t.Fatalf("Serialization error: %s", err.Error())
    }

    msg2, err := Deserialize(&buf)
    if err != nil {
        t.Fatalf("Deserialization error: %s", err.Error())
    }

    t.Assert(msg2.GetName(), Equals, msg.GetName())
    t.Assert(len(msg2.Headers), Equals, len(msg.Headers))
    t.Assert(msg2.Header("hdr1"), Equals, msg.Header("hdr1"))
    t.Assert(msg2.Header("hdr2"), Equals, msg.Header("hdr2"))
    t.Assert(len(msg2.BodyParts), Equals, len(msg.BodyParts))
    t.Assert(msg2.BodyPart("bp1"), DeepEquals, msg.BodyPart("bp1"))
    t.Assert(msg2.BodyPart("bp2"), DeepEquals, msg.BodyPart("bp2"))
}
