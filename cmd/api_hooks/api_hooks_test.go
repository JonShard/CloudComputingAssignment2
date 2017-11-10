package api_hooks

import "testing"

//Tests that the correct values are returned from the database.
func TestGetLatestCurrency(t *testing.T) {

	testVal := GetLatestCurrency("EUR", "EUR")
	if testVal != 1 {
		t.Error("\nFailed to get currencies EUR to EUR:\n\tShould had been 1.\n\tIt is ", testVal)
	}

}
