package p2p

import (
	"bytes"

	xdr "github.com/stellar/go-xdr/xdr3"
)

type Message struct {
	V  int32
	V0 MessageV0
}

func (m *Message) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	_, err := xdr.Marshal(&b, m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (m *Message) UnmarshalBinary(b []byte) error {
	_, err := xdr.Unmarshal(bytes.NewReader(b), m)
	if err != nil {
		return err
	}
	return nil
}

type MessageV0 struct {
	Chain      Chain
	Body       []byte
	Signatures [][]byte
}

func (m *MessageV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	_, err := xdr.Marshal(&b, m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (m *MessageV0) UnmarshalBinary(b []byte) error {
	_, err := xdr.Unmarshal(bytes.NewReader(b), m)
	if err != nil {
		return err
	}
	return nil
}

type Chain int32

func (c Chain) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	_, err := xdr.Marshal(&b, &c)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *Chain) UnmarshalBinary(b []byte) error {
	_, err := xdr.Unmarshal(bytes.NewReader(b), c)
	if err != nil {
		return err
	}
	return nil
}

const (
	ChainStellar  Chain = 0
	ChainEthereum Chain = 1
)
