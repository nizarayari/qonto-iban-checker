package main

import (
	"strings"
)

type Format struct {
	CountryCode    string
	Country        string
	IbanLength     int
	BankCodeLength int
}

type Bank struct {
	BankCode string
	Bic      string
	Name     string
}

func getBic(iban string) Bank {
	countryCode, Bban := iban[:2], iban[4:]
	countryCode = strings.ToUpper(countryCode)

	var resultFormat Format
	var resultBank Bank

	db.Debug().Table("formats_models").Where(&Format{CountryCode: countryCode}).Scan(&resultFormat)
	bankCode := Bban[:resultFormat.BankCodeLength]

	db.Debug().Table("banks_models").Where(&Bank{BankCode: bankCode}).Scan(&resultBank)

	return resultBank
}
