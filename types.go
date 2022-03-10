package postalcodes

import (
	"time"
)

type PostalCodeMap struct {
	BuiltAt        time.Time           `json:"built_at"`
	Postalcodes    map[string][]string `json:"postalcodes"`
	Municipalities map[string]string   `json:"municipalities"`
}

type Municipality struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// MunicipalityCodesOf retuns list of municipalities codes associated with a postal code
func (p *PostalCodeMap) MunicipalityCodesOf(code string) []string {
	m, found := p.Postalcodes[code]
	if found {
		return m
	}
	return nil
}

// MunicipalitiesOf returns list of Municipality with a postal code
func (p *PostalCodeMap) MunicipalitiesOf(code string) []Municipality {
	m, found := p.Postalcodes[code]
	if !found {
		return nil
	}
	r := make([]Municipality, 0, len(m))
	for _, code := range m {
		label, _ := p.Municipalities[code]
		r = append(r, Municipality{Code: code, Label: label})
	}
	return r
}
