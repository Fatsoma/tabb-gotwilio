package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// LookupType -Indicates the type of information you would like returned with your request
type LookupType string

// LookupTypeCarrier - Carrier Information LookupType
const LookupTypeCarrier LookupType = "carrier"

// LookupTypeCallerName - Caller Name Information LookupType
const LookupTypeCallerName LookupType = "caller-name"

// LookupOptions to be made alongside the lookup request
type LookupOptions struct {
	AddOns      string
	CountryCode string
	Type        LookupType
}

// CallerNameDetails -
type CallerNameDetails struct {
	CallerName string `json:"caller_name"`
	CallerType string `json:"caller_type"`
	ErrorCode  int    `json:"error_code"`
}

// CarrierDetails -
type CarrierDetails struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	ErrorCode         int    `json:"error_code"`
}

// FraudDetails -
type FraudDetails struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	AdvancedLineType  string `json:"advanced_line_type"`
	CallerName        string `json:"caller_name"`
	IsPorted          bool   `json:"is_ported"`
	LastPortedDate    string `json:"last_ported_date"`
	ErrorCode         int    `json:"error_code"`
}

// AddOnDetails -
type AddOnDetails struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Results interface{} `json:"results"`
}

// LookupResponse is returned after a lookup request is made to Twilio
type LookupResponse struct {
	CallerName     CallerNameDetails `json:"caller_name"`
	Carrier        CarrierDetails    `json:"carrier"`
	Fraud          FraudDetails      `json:"fraud"`
	AddOns         AddOnDetails      `json:"add_ons"`
	CountryCode    string            `json:"country_code"`
	NationalFormat string            `json:"national_format"`
	PhoneNumber    string            `json:"phone_number"`
	URL            string            `json:"url"`
}

// GetLookup uses twillio to get information about a phone number
// See https://www.twilio.com/docs/lookup/api for more information.
func (twilio *Twilio) GetLookup(phoneNumber string, lookupOptions LookupOptions) (lookupResponse *LookupResponse, exception *Exception, err error) {
	twilioUrl := twilio.LookupUrl + "/v1/PhoneNumbers/" + phoneNumber

	q := url.Values{}
	q.Add("AddOns", lookupOptions.AddOns)
	q.Add("CountryCode", lookupOptions.CountryCode)
	q.Add("Type", string(lookupOptions.Type))

	res, err := twilio.getWithParams(twilioUrl, q.Encode())

	if err != nil {
		return lookupResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return lookupResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)
		return lookupResponse, exception, err
	}

	lookupResponse = new(LookupResponse)
	err = json.Unmarshal(responseBody, lookupResponse)

	return lookupResponse, exception, err
}
