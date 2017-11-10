package api_hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JonShard/CloudComputingAssignment2/mongodb"
)

//InvokeAllHooks Invokes all the webhooks in the database if outside.
func InvokeAllHooks() {

	var hooksCollection = new(mongodb.MongoDB)
	hooksCollection.DatabaseURL = mongodb.MongoDatabaseURL
	hooksCollection.DatabaseName = mongodb.MongoCollectionDatabaseName
	hooksCollection.CollectionName = mongodb.HooksCollectionCollectionName
	hooksCollection.Init()

	hookEntries, succeeded := hooksCollection.GetAllWebookEntries()
	if !succeeded {
		fmt.Println(http.StatusText(500), "\nGetting the stuff messed up.")
	} else {

		length := len(hookEntries)

		for i := 1; i < length; i++ {

			latest := GetLatestCurrency(hookEntries[i].BaseCurrency, hookEntries[i].TargetCurrency)
			if latest != -1 {

				if latest < hookEntries[i].MinTriggerValue || latest > hookEntries[i].MaxTriggerValue {
					fmt.Println("\n\nInvoking a webhook.\nIs outside params: min=", hookEntries[i].MinTriggerValue, " max=", hookEntries[i].MaxTriggerValue, " val=", latest)

					tempURL := hookEntries[i].HookURL
					hookEntries[i].HookID = 0 // should be set to nil. So omitempty will work.
					hookEntries[i].HookURL = ""
					hookEntries[i].CurrentRate = latest

					body := new(bytes.Buffer)
					json.NewEncoder(body).Encode(hookEntries[i])
					resp, err := http.Post(tempURL, "application/json", body)
					if err != nil {
						fmt.Println("Invoke sent, response: ", resp)
					} else {
						fmt.Println("Failded: ", err)
					}
				}
			} else {
				fmt.Println("Invoking webhook failed because getting currencyEntry failed.")
			}
		}
	}

}

//StoreCurrencies gets the currencies from fixer.io, and stores them in the db.
func StoreCurrencies(base string) {

	var currenciesCollection = new(mongodb.MongoDB)
	currenciesCollection.DatabaseURL = mongodb.MongoDatabaseURL
	currenciesCollection.DatabaseName = mongodb.MongoCollectionDatabaseName
	currenciesCollection.CollectionName = mongodb.CurrenciesCollectionCollectionName
	currenciesCollection.Init()

	entry, err := GetFixerEntry(base)
	if err != nil {

		fmt.Println("Failed to get entry from Fixer.io")
	}
	currenciesCollection.AddCurrency(entry)

	fmt.Println("Well, it runs.(store)")

}

//GetLatestCurrency returns the most recent currency record in the db.
func GetLatestCurrency(base string, target string) float64 {

	if base == target { //A currency is worth 1 of itself.
		return 1
	}

	var currenciesCollection = new(mongodb.MongoDB)
	currenciesCollection.DatabaseURL = mongodb.MongoDatabaseURL
	currenciesCollection.DatabaseName = mongodb.MongoCollectionDatabaseName
	currenciesCollection.CollectionName = mongodb.CurrenciesCollectionCollectionName
	currenciesCollection.Init()

	var entry mongodb.CurrencyEntry
	var succeeded bool

	localTime := time.Now().Local().String()
	var parts []string
	parts = strings.Split(localTime, " ")
	localTime = parts[0]

	fmt.Println("API: ", localTime)

	entry, succeeded = currenciesCollection.GetCurrencyEntry(string(localTime))
	if succeeded {

		var targetValue float64
		if base == "EUR" {

			targetValue = entry.Rates[target] //Convert target(relative to EUR) to target(relative to base).
		} else {

			targetValue = entry.Rates[target] / entry.Rates[base] //Convert target(relative to EUR) to target(relative to base).
		}
		return targetValue
	}
	return -1
}

//GetFixerEntry gets the latest entry from fixer.io and returns it. Error if fails.
func GetFixerEntry(base string) (mongodb.CurrencyEntry, error) {

	var entry mongodb.CurrencyEntry

	site := "http://api.fixer.io/latest?base="

	currencyResponce, err := http.Get(site + base) //Get the http responce from fixer.
	if err != nil {
		return entry, err
	}

	json.NewDecoder(currencyResponce.Body).Decode(&entry) //Decode json into map.
	return entry, err
}
