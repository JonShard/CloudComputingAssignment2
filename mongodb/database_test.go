package mongodb

import (
	"testing"
)

func TestCountEntries(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	for i := 1; i <= 5; i++ {
		testEntry := createTestHook(i)
		testDB.AddHook(testEntry)
	}

	testVal := testDB.CountEntries()
	if testVal != 5 {
		t.Error("Failed to count hook entries.\n\tCount should had been 5\n\tIt is ", testVal)
	}
	testDB.DropCollection() //Must be in the last function in the test. Tears down the testDB for next time.
}

func TestAddHook(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	testEntryA := createTestHook(1)

	errA := testDB.AddHook(testEntryA)

	if errA != nil {
		t.Error("Failed to add new hook entry to testDB.\n\tError from function: ", errA)
	}

	testVal := testDB.CountEntries()
	if testVal != 1 {
		t.Error("Failed to add new hook entry to testDB.\n\tCount should had been 1\n\tIt is ", testVal, "\n\tError from function: ", errA)
	}

	testEntryB := createTestHook(1) //Tests that no duplicats PK constraint works.

	errB := testDB.AddHook(testEntryB)

	if errB == nil {
		t.Error("Failed to check that no duplicate PK constraint works.")
	}

	testDB.DropCollection() //Must be in the last function in the test. Tears down the testDB for next time.
}

func TestDropCollection(t *testing.T) {
	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	for i := 1; i <= 5; i++ {
		testEntry := createTestHook(i)
		testDB.AddHook(testEntry)
	}

	beforeCount := testDB.CountEntries()
	testDB.DropCollection()
	afterCount := testDB.CountEntries()

	if afterCount != 0 {
		t.Error("Failed to delete all webhook entries from testDB.\n\tcount before was", beforeCount, ", after delete should had been 0\n\tIt is ", afterCount)
	}
}

func TestAddCurrency(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = CurrenciesTestCollectionCollectionName
	testDB.Init()

	testEntry := createTestCurrency()

	err := testDB.AddCurrency(testEntry)

	if err != nil {
		t.Error("Failed to add new currency entry to testDB.\n\tError from function: ", err)
	}

	testVal := testDB.CountEntries()
	if testVal != 1 {
		t.Error("Failed to add new currency entry to testDB.\n\tCount should had been 1\n\tIt is ", testVal, "\n\tError from function: ", err)
	}

	testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func TestGetCurrencyEntry(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = CurrenciesTestCollectionCollectionName
	testDB.Init()

	testEntryA := createTestCurrency()
	testEntryA.Date = "0000-00-01"
	testEntryA.Base = "Test"

	errA := testDB.AddCurrency(testEntryA)

	if errA != nil {
		t.Error("Failed to add currency for testing", errA)
	}

	testEntryB, succeeded := testDB.GetCurrencyEntry("0000-00-01")

	if !succeeded {
		t.Error("Failed getting the entry form the database(probably does not exist)")
	}

	if testEntryB.Base != "Test" {
		t.Error("Failed to retrieve data stored in entry from db.")
	}

	testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func TestGetWebhookEntry(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	testEntryA := createTestHook(1)
	testEntryA.BaseCurrency = "Test"

	errA := testDB.AddHook(testEntryA)

	if errA != nil {
		t.Error("Failed to add hook for testing", errA)
	}

	testEntryB, succeeded := testDB.GetWebhookEntry(1)

	if !succeeded {
		t.Error("Failed getting the entry form the database(probably does not exist)")
	}

	if testEntryB.BaseCurrency != "Test" {
		t.Error("Failed to retrieve data stored in entry from db.")
	}

	testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func TestGetAllWebhookEntries(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	testEntryA := createTestHook(1)
	testEntryB := createTestHook(2)

	errA := testDB.AddHook(testEntryA)
	errB := testDB.AddHook(testEntryB)

	if errA != nil {
		t.Error("Failed to add new hook entry to testDB.\n\tError from function: ", errA)
	}
	if errB != nil {
		t.Error("Failed to add new hook entry to testDB.\n\tError from function: ", errB)
	}

	testEntries, succeeded := testDB.GetAllWebookEntries()
	if !succeeded {
		t.Error("Failed getting the entries form the database(empty database)")
	}

	testVal := len(testEntries)
	if testVal != 2 {
		t.Error("Failed to get all hook entries from testDB.\n\tLength of returned array should had been 2\n\tIt is ", testVal)
	}
	testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func TestGetAllCurrencyEntries(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = CurrenciesTestCollectionCollectionName
	testDB.Init()

	testEntryA := createTestCurrency()
	testEntryA.Date = "0000-00-01"
	testEntryB := createTestCurrency()
	testEntryB.Date = "0000-00-02"

	errA := testDB.AddCurrency(testEntryA)
	errB := testDB.AddCurrency(testEntryB)

	if errA != nil {
		t.Error("Failed to add new currency entry to testDB.\n\tError from function: ", errA)
	}
	if errB != nil {
		t.Error("Failed to add new currency entry to testDB.\n\tError from function: ", errB)
	}

	testEntries, succeeded := testDB.GetAllCurrencyEntries()
	if !succeeded {
		t.Error("Failed getting the entries form the database(empty database)")
	}

	testVal := len(testEntries)
	if testVal != 2 {
		t.Error("Failed to get all currency entries from testDB.\n\tLength of returned array should had been 2\n\tIt is ", testVal)
	}
	testDB.DropCollection() //Must be as the last function in the test. Tears down the testDB for next time.
}

func TestDeleteWebhook(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = MongoDatabaseURL
	testDB.DatabaseName = MongoCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	errB := testDB.DeleteWebhook(-1)
	if errB == nil {
		t.Error("Failed; did not get error when trying to delete webhook entry from testDB that does not exist:", errB)
	}

	testEntryA := createTestHook(1)
	errA := testDB.AddHook(testEntryA)

	if errA != nil {
		t.Error("Failed to add new hook entry for testing.\n\tError from function: ", errA)
	} else {

		err := testDB.DeleteWebhook(1)
		if err != nil {
			t.Error("Failed to delete webhook entry from testDB.", err)
		}

		testVal := testDB.CountEntries()
		if testVal != 0 {
			t.Error("Failed to delete webhook entry from testDB.\n\tcount after delete should had been 0\n\tIt is ", testVal)
		}
	}
}

func createTestHook(id int) WebhookEntry {
	var testEntry WebhookEntry
	testEntry.HookID = id
	testEntry.HookURL = "testURL"
	testEntry.BaseCurrency = "EUR"
	testEntry.TargetCurrency = "NOK"
	testEntry.MinTriggerValue = 9.1
	testEntry.MaxTriggerValue = 14.3
	return testEntry
}

func createTestCurrency() CurrencyEntry {

	var testEntry CurrencyEntry

	testEntry.Base = "EUR"
	testEntry.Date = "5555-55-55"

	ratesMap := map[string]float64{
		"AUD": 1.5382,
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
