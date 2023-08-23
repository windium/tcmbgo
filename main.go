package github.com/windium/tcmbgo

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// XML Struct
type tarihDate struct {
	XMLname  xml.Name `xml:"Tarih_Date"`
	Date     string   `xml:"Date,attr"`
	Currency []currency
}

type currency struct {
	Kod             string `xml:"Kod,attr"`
	CurrencyCode    string `xml:"CurrencyCode,attr"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`
	CrossRateUSD    string `xml:"CrossRateUSD"`
	CrossRateOther  string `xml:"CrossRateOther"`
}

type ExchangeRates struct {
	Date     string
	Currency []Currencies
}

type Currencies struct {
	CurrencyCode    string
	Name            string
	ForexBuying     float64
	ForexSelling    float64
	BanknoteBuying  float64
	BanknoteSelling float64
	CrossRateUSD    float64
	CrossRateOther  float64
}

const url = "https://www.tcmb.gov.tr/kurlar/"

// Get data from TCMB using timestamp.
func GetData(timestamp int64) *ExchangeRates {
	date := time.Unix(timestamp, 0)
	fDate := date.Format("2006.01.02")
	dateS := strings.Split(fDate, ".")

	fUrl := url + dateS[0] + dateS[1] + "/" + dateS[2] + dateS[1] + dateS[0] + ".xml"

	res, err := http.Get(fUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed code: %d", res.StatusCode)
	}
	data := parseBody(body)
	return data
}

func parseBody(body []byte) *ExchangeRates {
	var tarihDate tarihDate

	xml.Unmarshal(body, &tarihDate)

	format := new(ExchangeRates)
	format.Date = tarihDate.Date
	format.Currency = make([]Currencies, len(tarihDate.Currency))
	for i, c := range tarihDate.Currency {
		format.Currency[i].CurrencyCode = c.CurrencyCode
		format.Currency[i].Name = c.CurrencyName
		format.Currency[i].ForexBuying, _ = strconv.ParseFloat(c.ForexBuying, 64)
		format.Currency[i].ForexSelling, _ = strconv.ParseFloat(c.ForexSelling, 64)
		format.Currency[i].BanknoteBuying, _ = strconv.ParseFloat(c.BanknoteBuying, 64)
		format.Currency[i].BanknoteSelling, _ = strconv.ParseFloat(c.BanknoteSelling, 64)
		format.Currency[i].CrossRateUSD, _ = strconv.ParseFloat(c.CrossRateUSD, 64)
		format.Currency[i].CrossRateOther, _ = strconv.ParseFloat(c.CrossRateOther, 64)
	}

	return format
}
