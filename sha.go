package main

import "encoding/hex"

type sha string

func (s sha) bytes() ([]byte, error) {
	return hex.DecodeString(string(s))
}

func (s sha) match(other []byte) bool {
	me, err := s.bytes()
	if err != nil {
		return false
	}
	if len(me) != len(other) {
		return false
	}
	for ix := range me {
		if me[ix] != other[ix] {
			return false
		}
	}
	return true
}
