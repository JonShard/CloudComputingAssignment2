package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/JonShard/CloudComputingAssignment2/cmd/api_hooks"
	"github.com/JonShard/CloudComputingAssignment2/mongodb"
)

// ExtractID returnes the ID parsed from the givern URL.
func ExtractID(URL string) (int, error) {

	var parts []string
	parts = strings.Split(URL, "/") //Extract ID from URL.
	unparcedID := parts[2]
	id, err := strconv.Atoi(unparcedID) // Parse ID to int.

	return id, err

}

//Will handle the registration of new webhooks. And access to existing hooks.
func webhookHandler(w http.ResponseWriter, r *http.Request) {

	var hooksCollection = new(mongodb.MongoDB)
	hooksCollection.DatabaseURL = mongodb.MongoDatabaseURL
	hooksCollection.DatabaseName = mongodb.MongoCollectionDatabaseName
	hooksCollection.CollectionName = mongodb.HooksCollectionCollectionName
	hooksCollection.Init()

	switch r.Method {
	case "POST":
		var registrationInterface map[string]interface{} //Will temporarly store the values for registration.
		json.NewDecoder(r.Body).Decode(&registrationInterface)

		url, errURL := url.ParseRequestURI(registrationInterface["webHookURL"].(string))
		if errURL != nil {
			fmt.Fprintln(w, http.StatusText(400)) //Bad request if fails.
			return
		}
		var entry = new(mongodb.WebhookEntry) // Will store the values for registration.
		entry.HookURL = url.Path
		fmt.Println("parsed URL: ", url.Path)
		entry.BaseCurrency = registrationInterface["baseCurrency"].(string)
		entry.TargetCurrency = registrationInterface["targetCurrency"].(string)
		entry.MinTriggerValue = registrationInterface["minTriggerValue"].(float64)
		entry.MaxTriggerValue = registrationInterface["maxTriggerValue"].(float64)

		pk := hooksCollection.CountEntries() + 1
		if pk == -1 {
			fmt.Fprintf(w, http.StatusText(500))
			return
		}

		fmt.Printf("" + string(pk)) //Print next primary key value for testing.
		entry.HookID = pk
		err := hooksCollection.AddHook(*entry)
		if err != nil {
			fmt.Fprintf(w, http.StatusText(500))
			return
		}

		fmt.Fprintln(w, pk)

		break

	case "GET":

		id, err := ExtractID(r.URL.Path)
		if err != nil {
			fmt.Fprintf(w, http.StatusText(400)) //Bad request if fails.
			return
		}

		var entry mongodb.WebhookExportEntry //Get entry from mongoDB:
		var succeeded bool
		entry, succeeded = hooksCollection.GetWebhookEntry(id)
		if succeeded != true {
			fmt.Fprintf(w, http.StatusText(500)) //If mongoDB fails, Internal server error.
			return
		}

		b := new(bytes.Buffer)           //Creates a fake Body to encode into. Works mostly like a string.
		json.NewEncoder(b).Encode(entry) //Encodes the struct to json so it can be sent.

		fmt.Fprintln(w, b) //Prints json to webpage.

		break
	case "DELETE":

		id, err := ExtractID(r.URL.Path)
		if err != nil {
			fmt.Fprintf(w, http.StatusText(400)) //Bad request if fails.
			return
		}

		err2 := hooksCollection.DeleteWebhook(id)
		if err2 != nil {
			fmt.Fprintln(w, http.StatusText(404), "\nNo registered webhook with id:", id) //If mongoDB fails, Bad requset. Not found
			return
		}
		break

	}

}

//Will return the latest value of the target currency based on base-currnecy.
func latestHandler(w http.ResponseWriter, r *http.Request) {

	var inputInterface map[string]interface{}
	json.NewDecoder(r.Body).Decode(&inputInterface)

	base := inputInterface["baseCurrency"].(string)
	target := inputInterface["targetCurrency"].(string)

	latest := api_hooks.GetLatestCurrency(base, target)
	if latest == -1 {

		fmt.Fprintln(w, http.StatusText(500), ": No entries for that date.")

	} else {

		fmt.Fprintln(w, "", latest)
	}

}

//Will return the average value of the target currency based on base-currnecy over last 3days.
func averageHandler(w http.ResponseWriter, r *http.Request) {

	var inputInterface map[string]interface{}
	json.NewDecoder(r.Body).Decode(&inputInterface)

	base := inputInterface["baseCurrency"].(string)
	target := inputInterface["targetCurrency"].(string)

	var currenciesCollection = new(mongodb.MongoDB)
	currenciesCollection.DatabaseURL = mongodb.MongoDatabaseURL
	currenciesCollection.DatabaseName = mongodb.MongoCollectionDatabaseName
	currenciesCollection.CollectionName = mongodb.CurrenciesCollectionCollectionName
	currenciesCollection.Init()

	entry, succeeded := currenciesCollection.GetAllCurrencyEntries()
	if !succeeded {
		//error
		fmt.Println("Getting the stuff messed up.")

	} else {

		length := len(entry)

		var targetValue float64
		if base == "EUR" {

			targetValue = (entry[length-1].Rates[target] +
				entry[length-2].Rates[target] +
				entry[length-3].Rates[target]) / 3 //Convert target(relative to EUR) to target(relative to base).

		} else {

			targetValue = (((entry[length-1].Rates[target] + entry[length-2].Rates[target] + entry[length-3].Rates[target]) / 3) / ((entry[length-1].Rates[base] + entry[length-2].Rates[base] + entry[length-3].Rates[base]) / 3)) //Convert target(relative to EUR) to target(relative to base).

		}

		fmt.Fprintln(w, "", targetValue)
	}
}

func evaluationHandler(w http.ResponseWriter, r *http.Request) {

	api_hooks.InvokeAllHooks()
}

func main() {

	http.HandleFunc("/exchange/", webhookHandler)
	http.HandleFunc("/exchange/latest/", latestHandler)
	http.HandleFunc("/exchange/average/", averageHandler)
	http.HandleFunc("/exchange/evaluationtrigger/", evaluationHandler)

	port := os.Getenv("PORT")
	http.ListenAndServe("127.0.0.1:"+port, nil) // Keep serving all requests that is recieved.
}
