package api_hooks

import (
	"testing"

	"github.com/JonShard/CloudComputingAssignment2/mongodb"
)

//testDBURL holds the url to the datebase.
const HooksTestCollectionURL = "0.0.0.0:27017"

//testDBDatabaseName holds the name of the database.
const HooksTestCollectionDatabaseName = "local"

//testDBCollectionName holds the name of the collection.
const HooksTestCollectionCollectionName = "Testwebhooks"

//CurrenciesCollectionURL holds the url to the datebase.
const CurrenciesTestCollectionURL = "0.0.0.0:27017"

//CurrenciesCollectionDatabaseName holds the name of the database.
const CurrenciesTestCollectionDatabaseName = "local"

//CurrenciesCollectionDatabaseName holds the name of the database.
const CurrenciesTestCollectionCollectionName = "Testcurrencies"

//Tests that the correct values are returned from the database.
func TestGetLatestCurrency(t *testing.T) {

	testVal := GetLatestCurrency("EUR", "EUR")
	if testVal != 1 {
		t.Error("\nFailed to get currencies EUR to EUR:\n\tShould had been 1.\n\tIt is ", testVal)
	}

	// var testDB = new(mongodb.MongoDB)
	// testDB.DatabaseURL = CurrenciesTestCollectionURL
	// testDB.DatabaseName = CurrenciesTestCollectionDatabaseName
	// testDB.CollectionName = CurrenciesTestCollectionCollectionName
	// testDB.Init()
	//
	// localTime := time.Now().Local().String()
	// var parts []string
	// parts = strings.Split(localTime, " ")
	// localTime = parts[0]
	//
	// testEntryA := createTestCurrency()
	// testEntryA.Date = "0000-00-00"
	// testEntryA.Rates["AUD"] = 888.88888
	//
	// testEntryB := createTestCurrency()
	// testEntryB.Date = localTime
	// fmt.Println("TEST: ", localTime)
	// testEntryA.Rates["AUD"] = 999.99999
	//
	// errA := testDB.AddCurrency(testEntryA)
	// errB := testDB.AddCurrency(testEntryB)
	//
	// if errA != nil {
	// 	t.Error("Failed to add new currency entry to testDB.\n\tError from function: ", errA)
	// }
	// if errB != nil {
	// 	t.Error("Failed to add new currency entry to testDB.\n\tError from function: ", errB)
	// }
	//
	// testVal = GetLatestCurrency("EUR", "AUD")
	//
	// if testVal == -1 {
	// 	t.Error("Failed to get currency entry from testDB. No entries for that date.")
	// }
	//
	// if testVal != 999.99999 {
	// 	t.Error("Failed to add new currency entry to testDB.")
	//
	// }

	//testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func createTestCurrency() mongodb.CurrencyEntry {

	var testEntry mongodb.CurrencyEntry

	testEntry.Base = "EUR"
	testEntry.Date = "5555-55-55"

	ratesMap := map[string]float64{
		"AUD": 999.99999,
		"BGN": 1.9558,
		"BRL": 3.6134,
		"CAD": 1.4963,
		"CHF": 1.169,
		"CNY": 7.8317,
		"CZK": 25.589,
		"DKK": 7.4429,
		"GBP": 0.88883,
		"HKD": 9.1971,
		"HRK": 7.5148,
		"HUF": 310.165,
		"IDR": 15995.0,
		"ILS": 41.1328,
		"INR": 16.514,
		"JPY": 134.41,
		"KRW": 1325.4,
		"MXN": 42.603,
		"MYR": 4.9915,
		"NOK": 9.4383,
		"NZD": 31.7123,
		"PHP": 61.039,
		"PLN": 43.2355,
		"RON": 4.5975,
		"RUB": 67.864,
		"SEK": 9.6858,
		"SGD": 1.6047,
		"THB": 39.103,
		"TRY": 4.397,
		"USD": 1.1785,
		"ZAR": 16.297,
	}

	testEntry.Rates = ratesMap
	return testEntry
}
