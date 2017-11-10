package mongodb

import (
	"testing"
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

func TestAddHook(t *testing.T) {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = HooksTestCollectionURL
	testDB.DatabaseName = HooksTestCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	var testEntry WebhookEntry
	testEntry.HookID = testDB.CountEntries() + 1
	testEntry.HookURL = "testURL"
	testEntry.BaseCurrency = "EUR"
	testEntry.TargetCurrency = "NOK"
	testEntry.MinTriggerValue = 9.1
	testEntry.MaxTriggerValue = 14.3

	err := testDB.AddHook(testEntry)

	if err != nil {
		t.Error("Failed to add new hook entry to testDB.\n\tError from function: ", err)
	}

	testVal := testDB.CountEntries()
	if testVal != 1 {
		t.Error("Failed to add new entry to testDB.\n\tCount should had been 1\n\tIt is ", testVal, "\n\tError from function: ", err)
	}
}

//Tears down then test webhook collection.
func tearDownHooks() {

	var testDB = new(MongoDB)
	testDB.DatabaseURL = HooksTestCollectionURL
	testDB.DatabaseName = HooksTestCollectionDatabaseName
	testDB.CollectionName = HooksTestCollectionCollectionName
	testDB.Init()

	testDB.DropWebhook()
}
