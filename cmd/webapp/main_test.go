package main

import "testing"

//Tests that the correct id is parsed.
func TestExtractID(t *testing.T) {

	testVal, err := ExtractID("127.0.0.1:8080/exchange/5")

	if testVal != 5 {
		t.Error("Failed to parse ID from URL.\n\tShould had been 5\n\tIt is ", testVal, "\n\tError from function: ", err)
	}

}
