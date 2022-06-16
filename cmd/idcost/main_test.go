package main

import "testing"

func Test_getBytes(t *testing.T) {
	result := getBytes("0x")

	if len(result) != 0 {
		t.Errorf("result: %d", len(result))
	}

	result = getBytes("0x01")

	if len(result) != 1 {
		t.Errorf("result: %d", len(result))
	}

	if result[0] != "01" {
		t.Errorf("got: %s, want 01", result[0])
	}

	result = getBytes("0x1234")

	if len(result) != 2 {
		t.Errorf("result: %d", len(result))
	}

	if result[0] != "12" || result[1] != "34" {
		t.Errorf("got: %s, want 01", result[0])
	}
}

func Test_getCost(t *testing.T) {

	result := getCost(make([]string, 0))

	if result != 0 {
		t.Errorf("Cost for no bytes hould be zero")
	}

	result = getCost([]string{"00"})

	if result != 4 {
		t.Errorf("Cost of zero byte should be 4, got %d", result)
	}

	result = getCost([]string{"01"})

	if result != 68 {
		t.Errorf("Cost of non-zero byte should be 68, got %d", result)
	}

	result = getCost([]string{"00", "01"})

	if result != 72 {
		t.Errorf("Cost of 0001 byte should be 72, got %d", result)
	}

}
