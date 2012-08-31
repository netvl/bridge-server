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
    NAME_SIZE_MAX           = 255   // Maximum size of message name
    HEADERS_NUMBER_MAX      = 255   // Maximum number of headers
    BODY_PARTS_NUMBER_MAX   = 255   // Maximum number of body parts
    HEADER_NAME_SIZE_MAX    = 255   // Maximum size of header name
    HEADER_VALUE_SIZE_MAX   = 65535 // Maximum size of header value
    BODY_PART_NAME_SIZE_MAX = 255   // Maximum size of body part name
)

// Represents serialization error, implements error interface.
type SerializeError struct {
    msg   string
    cause error
}

// Creates new serialization error with given cause and message, which is
// formatted using fmt.Sprintf and additional arguments.
func makeErrorf(cause error, msg string, args ...interface{}) *SerializeError {
    return &SerializeError{fmt.Sprintf(msg, args...), cause}
}

// Creates new serialization error with given cause and message.
func makeError(cause error, msg string) *SerializeError {
    return &SerializeError{msg, cause}
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
    size   uint32
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

// Message represents a message which bridge nodes send to each other. Message has a name,
// a set of named headers and a set of named body parts.
type Message struct {
    name string
    headers map[string]string
    bodyParts map[string]*BodyPart
}

func CreateEmpty() *Message {
    return CreateWithName("")
}

// Create creates a message with given name and empty headers and body.
func CreateWithName(name string) *Message {
    return &Message{name, make(map[string]string), make(map[string]*BodyPart)}
}

// GetName returns message name.
func (this *Message) GetName() string {
    return this.name
}

// GetHeader returns a header by its name.
func (this *Message) GetHeader(name string) string {
    return this.headers[name]
}

// GetHeaderNames returns a slice containing all header names in lexicographical order.
func (this *Message) GetHeaderNames() []string {
    names := make([]string,len(this.headers))
    i := 0
    for k := range this.headers {
        names[i] = k
        i++
    }
    sort.Strings(names)
    return names
}

// ForEachHeader executes certain action for each header, providing it with header name
// and value.
func (this *Message) ForEachHeader(f func (string, string)) {
    for k, v := range this.headers {
        f(k, v)
    }
}

// GetBodyPart returns body part by its name.
func (this *Message) GetBodyPart(name string) *BodyPart {
    return this.bodyParts[name]
}

// GetBodyPartNames returns a slice containing all body part names in lexicographical order.
func (this *Message) GetBodyPartNames() []string {
    names := make([]string,len(this.bodyParts))
    i := 0
    for k := range this.bodyParts {
        names[i] = k
        i++
    }
    sort.Strings(names)
    return names
}

// ForEachBodyPart executes certain action for each body part, providing it with body part name
// and content.
func (this *Message) ForEachBodyPart(f func (string, *BodyPart)) {
    for k, v := range this.bodyParts {
        f(k, v)
    }
}

// SetName sets the message name.
func (this *Message) SetName(name string) {
    this.name = name
}

// SetHeader sets message header value of the header with given name.
func (this *Message) SetHeader(name, value string) {
    this.headers[name] = value
}

// SetBodyPart sets message body part content of the body part with given name.
func (this *Message) SetBodyPart(name string, value *BodyPart) {
    this.bodyParts[name] = value
}

// Message format
// _, __, ____ stand for size fields of corresponding number of bytes; multibyte sizes are little-endian
// {..} stands for field value
// Spaces, newlines and '--'-started comments should be ignored

// MSG____                    -- Message mark and 4-byte size
// _{name}                    -- Message name size and the name itself
// _                          -- Number of headers
// _{h1-name}_{h1-value}      -- Header name size, header name, header value size, header value
// ...
// _{hn-name}_{hn-value}
// _                          -- Number of body parts
// _{bp1-name}____{bp1-body}  -- Body part name size, body part name, body part size, body part content
// ...
// _{bpm-name}____{bpm-body}

func computeHeadersSize(msg *Message) uint32 {
    var headersLength uint32 = 1  // Headers number
    for k, v := range msg.headers {
        // Header name size + header name + header value size + header value
        headersLength += 1 + uint32(len(k)) + 2 + uint32(len(v))
    }
    return headersLength
}

func computeBodyPartsSize(msg *Message) uint32 {
    var bodyPartsLength uint32 = 1 // Body parts number
    for k, v := range msg.bodyParts {
        // Body part name size + body part name + body part size + body part
        bodyPartsLength += 1 + uint32(len(k)) + 4 + uint32(v.size)
    }
    return bodyPartsLength
}

func computeMessageSize(msg *Message) uint32 {
    var totalLength uint32 = 0  // 'MSG' initial string plus message size field are not considered as size

    // Name size field + name
    nameLength := uint32(len(msg.name)) + 1
    totalLength += nameLength

    totalLength += computeHeadersSize(msg)
    totalLength += computeBodyPartsSize(msg)

    return uint32(totalLength)
}

func serializeBodyPart(name string, body *BodyPart, dst io.Writer) error {
    // Write body part name and size
    buf := make([]byte,1 + len(name) + 4)
    buf[0] = uint8(len(name))
    copy(buf[1:], []byte(name))
    binary.LittleEndian.PutUint32(buf[len(name) + 1:], body.size)
    if _, err := dst.Write(buf); err != nil {
        return makeError(err, "Error writing body part name and size")
    }

    // Copy exactly body.size bytes from body part reader to the output
    if _, err := io.CopyN(dst, body.reader, int64(body.size)); err != nil {
        return makeError(err, "Error copying body part to output stream")
    }

    return nil
}

func validateMessage(msg *Message) error {
    if len(msg.name) > NAME_SIZE_MAX {
        return makeErrorf(nil, "Message size %d exceeds maximum of %d bytes",
            len(msg.name), NAME_SIZE_MAX)
    }

    if len(msg.headers) > HEADERS_NUMBER_MAX {
        return makeErrorf(nil, "Headers number %d exceeds maximum of %d records",
            len(msg.headers), HEADERS_NUMBER_MAX)
    }

    if len(msg.bodyParts) > BODY_PARTS_NUMBER_MAX {
        return makeErrorf(nil, "Body parts number %d exceeds maximum of %d records",
            len(msg.bodyParts), BODY_PARTS_NUMBER_MAX)
    }

    for k, v := range msg.headers {
        if len(k) > HEADER_NAME_SIZE_MAX {
            return makeErrorf(nil, "Message header '%s' name size %d exceeds maximum of %d bytes",
                k, len(k), HEADER_NAME_SIZE_MAX)
        }
        if len(v) > HEADER_VALUE_SIZE_MAX {
            return makeErrorf(nil, "Message header '%s' value size %d exceeds maximum of %d bytes",
                k, len(v), HEADER_VALUE_SIZE_MAX)
        }
    }

    for k, _ := range msg.bodyParts {
        if len(k) > BODY_PART_NAME_SIZE_MAX {
            return makeErrorf(nil, "Message body part '%s' name size %d exceeds maximum of %d bytes",
                k, len(k), BODY_PART_NAME_SIZE_MAX)
        }
    }

    return nil
}

// Serialize serializes a message to the provided io.Writer. Returns SerializeError instance in case of error.
func Serialize(msg *Message, dst io.Writer) error {
    // Validate message lengths
    if err := validateMessage(msg); err != nil {
        return err
    }

    // Write message mark
    if _, err := dst.Write([]byte("MSG")); err != nil {
        return makeError(err, "Error writing message mark")
    }

    // Write message size
    messageSize := computeMessageSize(msg)
    buf := make([]byte, 4)
    binary.LittleEndian.PutUint32(buf, messageSize)
    if _, err := dst.Write(buf); err != nil {
        return makeError(err, "Error writing message size")
    }

    // Write message name
    buf = make([]byte, 1 + len(msg.name))
    buf[0] = uint8(len(msg.name))
    copy(buf[1:], []byte(msg.name))
    if _, err := dst.Write(buf); err != nil {
        return makeError(err, "Error writing message name")
    }

    // Write message headers
    // Using GetHeaderNames here since we are required to have deterministic headers order
    // We are writing headers in one piece
    buf = make([]byte, computeHeadersSize(msg))
    i := 0
    buf[i] = uint8(len(msg.headers))
    i++
    for _, k := range msg.GetHeaderNames() {
        v := msg.headers[k]

        // Write header name
        buf[i] = uint8(len(k))
        i += 1
        copy(buf[i:i + len(k)], []byte(k))
        i += len(k)

        // Write header value
        binary.LittleEndian.PutUint16(buf[i:i + 2], uint16(len(v)))
        i += 2
        copy(buf[i:i + len(v)], []byte(v))
        i += len(v)
    }
    if _, err := dst.Write(buf); err != nil {
        return makeError(err, "Error writing headers")
    }

    // Write message body parts
    // Same argument as for GetHeaderNames
    bsz := [1]byte{uint8(len(msg.bodyParts))}
    if _, err := dst.Write(bsz[:]); err != nil {
        return makeError(err, "Error writing body parts number")
    }
    for _, k := range msg.GetBodyPartNames() {
        v := msg.bodyParts[k]
        if err := serializeBodyPart(k, v, dst); err != nil {
            return makeErrorf(err, "Error writing body part '%s'", k)
        }
    }

    return nil
}

func readUint8(src io.Reader) (uint8, error) {
    var s [1]byte

    if _, err := src.Read(s[:]); err != nil {
        return 0, makeError(err, "Error reading uint8");
    }

    return s[0], nil
}

func readUint16(src io.Reader) (uint16, error) {
    var s [2]byte
    ss := s[:]

    if n, err := src.Read(ss); n < 2 {
        return 0, makeError(nil, "Unexpected end of stream while reading uint16")
    } else if err != nil {
        return 0, makeError(err, "Error reading uint16")
    }

    return binary.LittleEndian.Uint16(ss), nil
}

func readUint32(src io.Reader) (uint32, error) {
    var s [4]byte
    ss := s[:]

    if n, err := src.Read(ss); n < 4 {
        return 0, makeError(nil, "Unexpected end of stream while reading uint32")
    } else if err != nil {
        return 0, makeError(err, "Error reading uint32")
    }

    return binary.LittleEndian.Uint32(ss), nil
}

// DeserializeMessageName starts deserializing a message, loading message name and size from the given io.Reader.
// Returns newly created message or error, if any.
func DeserializeMessageName(src io.Reader) (*Message, error) {
    var msg *Message = CreateEmpty()

    // Read and check message mark and message size
    buf := make([]byte, 7)
    n, err := src.Read(buf)
    if n < 7 {
        return nil, makeError(nil, "Unexpected end of stream while reading message mark and size")
    } else if err != nil {
        return nil, makeError(err, "Error reading message mark and size")
    } else if (string(buf[0:3]) != "MSG") {
        return nil, makeError(err, "Message header was not found")
    }

    // Read message name
    sz, err := readUint8(src)
    if err != nil {
        return nil, makeError(err, "Error reading message name size")
    }
    buf = make([]byte, int(sz))
    if n, err := src.Read(buf); n < int(sz) {
        return nil, makeError(err, "Unexpected end of stream while reading message name")
    } else if err != nil {
        return nil, makeError(err, "Error reading message name")
    }
    msg.name = string(buf)

    return msg, nil
}

// DeserializeMessageHeaders deserializes message headers from the given io.Reader and puts them
// into specified Message. This function is supposed to be called after DeserializeMessageName.
func DeserializeMessageHeaders(src io.Reader, msg *Message) error {
    // Read number of headers
    sz, err := readUint8(src)
    if err != nil {
        return makeError(err, "Error reading headers number")
    }

    // Read the headers themselves
    for i := 1; i <= int(sz); i++ {
        // Read header name
        hsz, err := readUint8(src)
        if err != nil {
            return makeErrorf(err, "Error reading %d header name size", i)
        }
        buf := make([]byte, int(hsz))
        if n, err := src.Read(buf); n < int(hsz) {
            return makeErrorf(err, "Unexpected end of stream while reading %d header name", i)
        } else if err != nil {
            return makeErrorf(err, "Error reading %d header name", i)
        }
        name := string(buf)

        // Read header value
        vsz, err := readUint16(src)
        if err != nil {
            return makeErrorf(err, "Error reading %d header name size", i)
        }
        buf = make([]byte, int(vsz))
        if n, err := src.Read(buf); n < int(vsz) {
            return makeErrorf(err, "Unexpected end of stream while reading %d header value", i)
        } else if err != nil {
            return makeErrorf(err, "Error reading %d header value", i)
        }
        value := string(buf)

        msg.headers[name] = value
    }

    return nil
}

// Represents a hook which can be installed into deserialization process.
type DeserializeHook func (name string, size uint32, src io.Reader) (bool, error)

// EmptyHook is a function of type DeserializeHook which always returns false;
// because of this deserializing body parts with this hook results in loading the message
// completely in memory
var EmptyHook = func (string, uint32, io.Reader) (bool, error) { return false, nil }

// Deserializes body parts from the given io.Reader and puts them into specified Message.
// This function is supposed to be called after DeserializeMessageHeader function.
//
// Is it possible to supply hook function which can do arbitrary work with body part.
//
// This function takes body part name, size and input stream and should return true if
// it has accepted this body part and processed it or false if simple body part processing
// is required. Simple processing means putting all body parts to byte array into memory.
//
// It is possible for this function to return non-nil error; then first return value
// is interpreted as follows: if it is true, then the processing should be stopped and an error
// should be signaled; if it is false, the processing should continue.
//
// It is expected that in case of success the hook will read exactly specified number of bytes.
func DeserializeMessageBodyParts(src io.Reader, msg *Message, hook DeserializeHook) error {
    // Read body parts number
    sz, err := readUint8(src)
    if err != nil {
        return makeError(err, "Error reading body parts number")
    }

    // Read body parts
    for i := 1; i <= int(sz); i++ {
        // Read body part name
        hsz, err := readUint8(src)
        if err != nil {
            return makeErrorf(err, "Error reading %d body part name size", i)
        }
        buf := make([]byte, int(hsz))
        if n, err := src.Read(buf); n < int(hsz) {
            return makeErrorf(err, "Unexpected end of stream while reading %d body part name", i)
        } else if err != nil {
            return makeErrorf(err, "Error reading %d body part name")
        }
        name := string(buf)

        // Read body part value
        vsz, err := readUint32(src)
        if err != nil {
            return makeErrorf(err, "Error reading %d body part content size", i)
        }
        // Try to apply a hook first
        r, err := hook(name, vsz, src)
        if err != nil {
            if r {
                return makeErrorf(err, "Processing of %d body part content failed", i)
            }
        } else if !r {
            // Load body part in memory
            buf := make([]byte, vsz)
            if n, err := src.Read(buf); n < int(vsz) {
                return makeErrorf(err, "Unexpected end of stream while reading %d body part name", i)
            } else if err != nil && !(err == io.EOF && i == int(sz)) {
                // The condition above filters the case when we got EOF error reading the last body part,
                // so it is not considered an error
                return makeErrorf(err, "Error reading %d body part name", i)
            }
            msg.bodyParts[name] = BodyPartFromSlice(buf)
        }
    }

    return nil
}
