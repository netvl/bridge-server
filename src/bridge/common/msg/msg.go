/**
 * Date: 24.08.12
 * Time: 16:29
 *
 * @author Vladimir Matveev
 */
package msg

import (
    "encoding/binary"
    "bytes"
    "io"
    "fmt"
    "sort"
)

const (
    // Size of buffer to use during serialize operations in bytes.
    BUFFER_SIZE = 1024
)

// Represents serialization error, implements error interface.
type SerializeError struct {
    msg string
    cause error
}

// Returns error message.
func (this *SerializeError) Error() string {
    return this.msg
}

// Returns cause of serialization error, which is another error instance.
func (this *SerializeError) Cause() error {
    return this.cause
}

// Represents single body part datum. It is backed by io.Reader since
// it is possible to have large chunks of data as body (e.g. file), so
// it is not wise to load it whole into memory.
type BodyPart struct {
    size uint32
    reader io.Reader
}

// Creates a body part from given slice of bytes.
func BodyPartFromSlice(data []byte) *BodyPart {
    return &BodyPart{uint32(len(data)), bytes.NewReader(data)}
}

// Creates a body part from given reader and its size.
func BodyPartFromReader(size uint32, reader io.Reader) *BodyPart {
    return &BodyPart{size, reader}
}

func (this *BodyPart) Size() uint32 {
    return this.size
}

func (this *BodyPart) Read(dst []byte) (int, error) {
    return this.reader.Read(dst)
}

// A message with which bridge nodes send to each other. Message has a name,
// a set of named headers and a set of named body parts.
type Message struct {
    name string
    headers map[string]string
    bodyParts map[string]*BodyPart
}

// Creates a message with given name and empty headers and body.
func Create(name string) *Message{
    return &Message{name, make(map[string]string), make(map[string]*BodyPart)}
}

// Returns message name.
func (this *Message) GetName() string {
    return this.name
}

// Returns a header by its name.
func (this *Message) GetHeader(name string) string {
    return this.headers[name]
}

// Returns a slice containing all header names in lexicographical order.
func (this *Message) GetHeaderNames() []string {
    names := make([]string, len(this.headers))
    i := 0
    for k := range this.headers {
        names[i] = k
        i++
    }
    sort.Strings(names)
    return names
}

// Executes certain action for each header, providing it with header name
// and value.
func (this *Message) ForEachHeader(f func (string, string)) {
    for k, v := range this.headers {
        f(k, v)
    }
}

// Returns body part by its name.
func (this *Message) GetBodyPart(name string) *BodyPart {
    return this.bodyParts[name]
}

// Returns a slice containing all body part names in lexicographical order.
func (this *Message) GetBodyPartNames() []string {
    names := make([]string, len(this.bodyParts))
    i := 0
    for k := range this.bodyParts {
        names[i] = k
        i++
    }
    sort.Strings(names)
    return names
}

// Executes certain action for each body part, providing it with body part name
// and content.
func (this *Message) ForEachBodyPart(f func (string, *BodyPart)) {
    for k, v := range this.bodyParts {
        f(k, v)
    }
}

// Sets message name.
func (this *Message) SetName(name string) {
    this.name = name
}

// Sets message header.
func (this *Message) SetHeader(name, value string) {
    this.headers[name] = value
}

// Sets message body part.
func (this *Message) SetBodyPart(name string, value *BodyPart) {
    this.bodyParts[name] = value
}

// Message format
// ____ stands for 4-byte size field
// MSG____{name}\0{h1-name}\0{h1-value}\0...{hn-name}\0{hn-value}\0\0{bp1-name}\0____{bp1-body}....

func computeMessageSize(message *Message) uint32 {
    totalLength := 0  // 'MSG' initial string plus size field are not considered as size

    nameLength := len(message.name) + 1
    totalLength += nameLength

    headersLength := 0
    for k, v := range message.headers {
        // Header name + zero char + header value + zero char
        headersLength += len(k) + 1 + len(v) + 1
    }
    headersLength += 1  // Terminating zero
    totalLength += headersLength

    bodyPartsLength := 0
    for k, v := range message.bodyParts {
        // Body part name + zero char + body part size + body part
        bodyPartsLength += len(k) + 1 + 4 + int(v.size)
    }
    totalLength += bodyPartsLength

    return uint32(totalLength)
}

func serializeBodyPart(name string, body *BodyPart, dst io.Writer) error {
    // Write body part name and size
    nameSizeBuf := make([]byte, len(name)+1+4)
    copy(nameSizeBuf, []byte(name))
    binary.LittleEndian.PutUint32(nameSizeBuf[len(name)+1:], body.size)
    if _, err := dst.Write(nameSizeBuf); err != nil {
        return &SerializeError{"Error writing body part name and size", err}
    }

    // Copy data from body part reader to the output
    buf := make([]byte, BUFFER_SIZE)
    for {
        n, err := body.reader.Read(buf)
        if err != nil {
            return &SerializeError{"Error reading body part from reader", err}
        }

        if _, err := dst.Write(buf[:n]); err != nil {
            return &SerializeError{"Error writing body part", err}
        }

        if n < BUFFER_SIZE {
            break;
        }
    }

    return nil
}

// Serializes message to provided io.Writer.
func Serialize(message *Message, dst io.Writer) error {
    // Write header
    if _, err := dst.Write([]byte("MSG")); err != nil {
        return &SerializeError{"Error writing message header", err}
    }

    // Write message size
    messageSize := computeMessageSize(message)
    sizeBuf := make([]byte, 4)
    binary.LittleEndian.PutUint32(sizeBuf, messageSize)
    if _, err := dst.Write(sizeBuf); err != nil {
        return &SerializeError{"Error writing message size", err}
    }

    // Write message name
    nameBuf := make([]byte, len(message.name) + 1)
    copy(nameBuf, []byte(message.name))
    if _, err := dst.Write(nameBuf); err != nil {
        return &SerializeError{"Error writing message name", err}
    }

    // Write message headers
    // Using GetHeaderNames since we are required to have deterministic headers order
    headersBuf := make([]byte, 0)
    for _, k := range message.GetHeaderNames() {
        v := message.headers[k]
        headersBuf = append(headersBuf, []byte(k)...)
        headersBuf = append(headersBuf, 0)
        headersBuf = append(headersBuf, []byte(v)...)
        headersBuf = append(headersBuf, 0)
    }
    headersBuf = append(headersBuf, 0)
    if _, err := dst.Write(headersBuf); err != nil {
        return &SerializeError{"Error writing headers", err}
    }

    // Write message body parts
    // Same argument for GetBodyPartNames
    for _, k := range message.GetBodyPartNames() {
        v := message.bodyParts[k]
        if err := serializeBodyPart(k, v, dst); err != nil {
            return &SerializeError{fmt.Sprintf("Error writing body part '%s'", k), err}
        }
    }

    return nil
}

// Deserializes a message from given io.Reader. Is it possible to supply hook
// function which can do arbitrary work with body part. This function takes
// body part name, size and input stream and should return true if it has accepted
// this body part and processed it or false if simple body part processing is required.
// It is possible for this function to return non-nil error, then first return value
// is interpreted as follows: if it is true, then the processing should be stopped and an error
// should be signaled; if it is false, the processing should continue. In first case
// the input reader will be fully exhausted first.
func Deserialize(reader io.Reader, hook func (string, uint32, io.Reader) (bool, error)) (*Message, error) {
    var msg Message

    headerSlice := make([]byte, 7)
    _, err := reader.Read(headerSlice)
    if err != nil {
        return nil, &SerializeError{"Error reading header", err}
    } else if (string(headerSlice[0:3]) != "MSG") {
        return nil, &SerializeError{"Message header was not found!", nil}
    }

//    messageSize := binary.LittleEndian.Uint32(headerSlice[3:7])

    return &msg, nil
}
