package mongodb

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

//HooksCollectionURL holds the url to the datebase.
const HooksCollectionURL = "0.0.0.0:27017"

//HooksCollectionDatabaseName holds the name of the database.
const HooksCollectionDatabaseName = "local"

//HooksCollectionCollectionName holds the name of the collection.
const HooksCollectionCollectionName = "webhooks"

//CurrenciesCollectionURL holds the url to the datebase.
const CurrenciesCollectionURL = "0.0.0.0:27017"

//CurrenciesCollectionDatabaseName holds the name of the database.
const CurrenciesCollectionDatabaseName = "local"

//CurrenciesCollectionDatabaseName holds the name of the database.
const CurrenciesCollectionCollectionName = "currencies"

// MongoDB contains the information required to access a collection.
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

//WebhookEntry contains the required information to invoke a webhook.
type WebhookEntry struct {
	HookID          int     `json:"hookID, omitempty"`
	HookURL         string  `json:"webhookURL, omitempty"`
	BaseCurrency    string  `josn:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	CurrentRate     float64 `json:"currentRate, omitempty"`
	MinTriggerValue float64 `json:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue"`
}

//WebhookExportEntry will be used when responding to GET request.
type WebhookExportEntry struct {
	HookURL         string  `json:"webhookURL, omitempty"`
	BaseCurrency    string  `josn:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	MinTriggerValue float64 `json:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue"`
}

//WebhookInvokeEntry will be used when invoking a webhook.
type WebhookInvokeEntry struct {
	BaseCurrency    string  `josn:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	CurrentRate     float64 `json:"currentRate, omitempty"`
	MinTriggerValue float64 `json:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue"`
}

//CurrencyEntry contains the data for the registered values for a day.
type CurrencyEntry struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

//Init Initialises the database.
func (db *MongoDB) Init() {

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var index mgo.Index
	if db.CollectionName == "webhooks" || db.CollectionName == "Testwebhooks" {

		index = mgo.Index{ //Constrains the collection to only have one entry with same ID.
			Key:        []string{"hookid"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}

	} else if db.CollectionName == "currencies" || db.CollectionName == "Testcurrencies" {
		index = mgo.Index{ //Constrains the collection to only have one entry with same ID.
			Key:        []string{"date"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}
	}
	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

}

func (db *MongoDB) AddHook(h WebhookEntry) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(h)

	return err
}

func (db *MongoDB) AddCurrency(c CurrencyEntry) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(c)

	return err
}

func (db *MongoDB) CountEntries() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()

	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

func (db *MongoDB) GetCurrencyEntry(selectedDate string) (CurrencyEntry, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var entry CurrencyEntry
	succeeded := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"date": selectedDate}).One(&entry)
	if err != nil {
		succeeded = false
	}

	return entry, succeeded
}

func (db *MongoDB) GetWebhookEntry(id int) (WebhookExportEntry, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var entry WebhookExportEntry
	succeeded := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"hookid": id}).One(&entry)
	if err != nil {
		succeeded = false
	}

	return entry, succeeded
}

func (db *MongoDB) GetAllCurrencyEntries() ([]CurrencyEntry, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var entry []CurrencyEntry
	succeeded := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(nil).All(&entry)
	if err != nil {
		succeeded = false
	}

	return entry, succeeded
}

func (db *MongoDB) GetAllWebookEntries() ([]WebhookEntry, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var entry []WebhookEntry
	succeeded := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(nil).All(&entry)
	if err != nil {
		succeeded = false
	}

	return entry, succeeded
}

func (db *MongoDB) DeleteWebhook(id int) error {

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Remove(bson.M{"hookid": id})

	if err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) DropCollection() {

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.DB(db.DatabaseName).C(db.CollectionName).RemoveAll(nil)

}
