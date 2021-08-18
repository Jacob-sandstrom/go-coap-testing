package coaptest

import (
	"bytes"
	"fmt"
	"net"
)

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

type Option struct {
	Number        byte
	Delta         byte
	Len           byte
	DeltaExtended []byte
	LenExtended   []byte
	Value         []byte
}

func ParseOption(payload []byte, byteIndex int, lastNumber byte) (Option, int) {
	o := Option{}
	o.Delta = payload[byteIndex] & 240 >> 4 // First four bits
	o.Len = payload[byteIndex] & 15         // Last four bits
	if o.Delta > 12 {
		// Todo: handle big delta
		panic("big delta not handled")
	}
	if o.Len > 12 {
		// Todo: handle big len
		panic("big len not handled")
	}
	o.Number = lastNumber + o.Delta
	o.Value = payload[byteIndex+1 : byteIndex+1+int(o.Len)]

	byteIndex += int(o.Len) + 1
	return o, byteIndex
}

func ParseCoap(payload []byte, len int) CoapMsg {
	msg := CoapMsg{}

	msg.Version = payload[0] & 192 >> 6 // First two bits
	msg.Type = payload[0] & 48 >> 4     // Bits three and four
	msg.TokenLen = payload[0] & 15      // Last four bits
	msg.Code = payload[1]
	msg.MsgID = payload[2:4]
	msg.Token = payload[4 : 4+msg.TokenLen]

	byteIndex := 4 + int(msg.TokenLen)
	var o Option
	for {
		if byteIndex == len {
			break
		}
		if payload[byteIndex] == 255 {
			byteIndex++
			break
		}
		o, byteIndex = ParseOption(payload, byteIndex, o.Number)
		msg.Options = append(msg.Options, o)
	}

	fmt.Printf("Parsed coap message: %+v\n", msg)
	return msg
}

func getCoapMsg(port int, adress string) CoapMsg {
	addr := net.UDPAddr{
		Port: port,
		// IP:   net.ParseIP(adress),
	}
	conn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var buf [1024]byte
	var msg CoapMsg
	for {
		rlen, _, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			panic(err.Error())
		}
		msg = ParseCoap(buf[:rlen], rlen)

		break
	}

	return msg
}

func compareOption(o0 Option, o1 Option) error {
	if o0.Number != o1.Number {
		return fmt.Errorf("Error wrong Version: Expected %v, Got %v", o0.Number, o1.Number)
	}
	if !bytes.Equal(o0.Value, o1.Value) {
		return fmt.Errorf("Error wrong Value: Expected %v, Got %v", o0.Value, o1.Value)
	}
	return nil
}

func containsOptions(expectedOptions []Option, options []Option) bool {
	for i, eo := range expectedOptions {
		fmt.Println(i)
		hasOption := false
		for _, o := range options {
			err := compareOption(eo, o)
			if err == nil {
				hasOption = true
				break
			}
		}

		if !hasOption {
			return false
		}
	}
	return true
}

func CompareCoap(expected CoapMsg, msg CoapMsg) error {
	if expected.Version != 0 {
		if expected.Version != msg.Version {
			return fmt.Errorf("Error wrong Version: Expected %v, Got %v", expected.Version, msg.Version)
		}
	}
	if expected.Type != 0 {
		if expected.Type != msg.Type {
			return fmt.Errorf("Error wrong Type: Expected %v, Got %v", expected.Type, msg.Type)
		}
	}
	if expected.TokenLen != 0 {
		if expected.TokenLen != msg.TokenLen {
			return fmt.Errorf("Error wrong TokenLen: Expected %v, Got %v", expected.TokenLen, msg.TokenLen)
		}
	}
	if expected.Code != 0 {
		if expected.Code != msg.Code {
			return fmt.Errorf("Error wrong Code: Expected %v, Got %v", expected.Code, msg.Code)
		}
	}
	if expected.MsgID != nil {
		if bytes.Compare(expected.MsgID, msg.MsgID) != 0 {
			return fmt.Errorf("Error wrong MsgID: Expected %v, Got %v", expected.MsgID, msg.MsgID)
		}
	}
	if expected.Token != nil {
		if bytes.Compare(expected.Token, msg.Token) != 0 {
			return fmt.Errorf("Error wrong Token: Expected %v, Got %v", expected.Token, msg.Token)
		}
	}
	if expected.Options != nil {
		if !containsOptions(expected.Options, msg.Options) {
			return fmt.Errorf("Error wrong Options: Expected %v, Got %v", expected.Options, msg.Options)

		}
	}
	if expected.Payload != nil {
		if bytes.Equal(expected.Payload, msg.Payload) {
			return fmt.Errorf("Error wrong Payload: Expected %v, Got %v", expected.Payload, msg.Payload)
		}
	}

	return nil
}

func Expect(expected CoapMsg) error {
	msg := getCoapMsg(2000, "127.0.0.1")

	return CompareCoap(expected, msg)
}
