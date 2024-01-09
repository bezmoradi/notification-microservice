package helpers

import "encoding/json"

func JsonValidator(messageBody string) (string, bool) {

	var jsonValue map[string]interface{}
	json.Unmarshal([]byte(messageBody), &jsonValue)

	dataValue := jsonValue["data"]
	emailBody, emailBodyIsValid := dataValue.(string)

	return emailBody, emailBodyIsValid
}
