# Usage

## Create the coap message you want
```go
type CoapMsg struct {
	Version  byte
	Type     byte
	TokenLen byte
	Code     byte
	MsgID    []byte
	Token    []byte
	Options  []Option
	Payload  []byte
}
```

## Call Expect(port, CoapMsg)

The Expect function starts a server at the given port and will process the first udp message recieved into a CoapMsg struct, comparing it to the given struct and returning an error if they don't match.

Only non nil fields in the CoapMsg are checked and excess options are ignored

## Example 
[Uses "github.com/plgd-dev/go-coap" as a client](https://github.com/plgd-dev/go-coap)

```go
package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	coaptest "github.com/Jacob-sandstrom/go-coap-testing"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
)

func SendCoap(port int) error {
	co, err := udp.Dial(fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return err
	}
	client := client.NewClient(co)
	client.Post(client.Context(), "/bs", 0, bytes.NewReader(nil))

	return nil
}

func TestSendCoap(t *testing.T) {
	testPort := 3000

	// Create the expected CoapMsg
	expected := coaptest.CoapMsg{Code: 2, Options: []coaptest.Option{{Number: 11, Value: []byte("bs")}}}

	go func() {
		time.Sleep(1 * time.Second) // Give the listener some time to start
		err := sendCoap(testPort)   // Call the function you want to test
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	err := coaptest.Expect(testPort, expected) // Start the listener
	if err != nil {
		t.Errorf(err.Error())
	}
}

```
