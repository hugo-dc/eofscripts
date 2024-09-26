package common

import (
	"encoding/hex"
	"testing"
)

func TestAddCode(t *testing.T) {
	eofObject := NewEOFObject()
	eofObject.AddCode([]byte{0xfe})

	expected := "ef000101000402000100010400000000800000fe"
	if eofObject.Code() != expected {
		t.Errorf("Error: Code is not parsed correctly")
		t.Errorf("Expected: \n%s", expected)
		t.Errorf("Result: \n%s", eofObject.Code())
	}
}

func TestParseEOF(t *testing.T) {
	eofCode := "ef000101000402000100010400000000800000fe"
	eofBytecode, err := hex.DecodeString(eofCode)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	result, err := ParseEOF(eofBytecode)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	expected := NewEOFObject()
	expected.Types = [][]int64{{0, 0x80, 0}}
	expected.CodeSections = [][]byte{{0xfe}}

	if !result.EqualTo(expected) {
		t.Errorf("Error: Code is not parsed correctly")
		t.Errorf("Expected: \n%s", expected.Describe())
		t.Errorf("Result: \n%s", result.Describe())
	}
}

func TestCalculateMaxStack(t *testing.T) {
	/*
		eofObject := NewEOFObject()

		eofObject.AddCode("60005e0100035b5b006014602760003960146000f3")
		eofObject.AddData("ef000101000402000100010300000000000000fe")

		eof_code := eofObject.Code(false, true)

		t.Errorf("Error: %s", eof_code)
	*/
}
