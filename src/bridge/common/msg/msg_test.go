/**
 * Date: 25.08.12
 * Time: 20:05
 *
 * @author Vladimir Matveev
 */
package msg

import (
    "testing"
    "bytes"
    "encoding/hex"
)

const (
    EXAMPLE1_HEXDUMP = "4d534750000000746573745f6d65737361676500686472310076616c75653100686472320076616c" +
        "756532000062703100040000000102030462703200160000005465737420626f6479207061727420636f6e74656e74"
)

func createMsg_Example1() *Message {
    msg := Create("test_message")
    msg.SetHeader("hdr1", "value1")
    msg.SetHeader("hdr2", "value2")
    msg.SetBodyPart("bp1", BodyPartFromSlice([]byte{1, 2, 3, 4}))
    msg.SetBodyPart("bp2", BodyPartFromSlice([]byte("Test body part content")))
    return msg
}

func TestMessageSerialization(t *testing.T) {
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
