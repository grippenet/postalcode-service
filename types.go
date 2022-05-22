// Types used by the server the builder
// The main structure is a registry containing Postal code to municipality entry and a registry of municipalities
// We only map the Postal code to Municipality
package postalcodes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Only 36000 municipalities for now, uint16
type MapIndex uint16

// unique create a unique list of index from a list of index
func unique(values []MapIndex) []MapIndex {
	set := make(map[MapIndex]struct{}, 0)
	var empty struct{}
	for _, v := range values {
		_, found := set[v]
		if !found {
			set[v] = empty
		}
	}
	keys := make([]MapIndex, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

// Municipality describes a known municipality with official code and a label
type Municipality struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// PostalCodeMap holds the registry of known municipalities and the postal code to municipalities index.
// A Postal code can be mapped to several municipalities
type PostalCodeMap struct {
	BuiltAt        time.Time                 `json:"built_at"`
	Postalcodes    map[string][]MapIndex     `json:"postalcodes"`
	Municipalities map[MapIndex]Municipality `json:"municipalities"`
	codes          map[string]MapIndex       // Index of municipalities by municipality code
}

// MapBuilder handles the creation of the municipalities registry, indexing and association to postal code
type MapBuilder struct {
	target  PostalCodeMap
	index   MapIndex
	indexes map[string]MapIndex
}

// LoadPostalCodeMap load the Postal code mapping registry from json
func LoadPostalCodeMap(file string) (*PostalCodeMap, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := PostalCodeMap{}

	err = json.Unmarshal(data, &p)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Build municipality code index index
	codes := make(map[string]MapIndex, len(p.Municipalities))
	for index, m := range p.Municipalities {
		codes[m.Code] = index
	}
	p.codes = codes

	return &p, nil
}

// MunicipalitiesOfPostal returns list of Municipality with a postal code
func (p *PostalCodeMap) MunicipalitiesOfPostal(postal string) []Municipality {
	m, found := p.Postalcodes[postal]
	if !found {
		return nil
	}
	r := make([]Municipality, 0, len(m))
	for _, index := range m {
		record := p.Municipalities[index]
		r = append(r, record)
	}
	return r
}

// LabelAt returns municipality label for a given municipality code
func (p *PostalCodeMap) LabelAt(code string) string {
	index, found := p.codes[code]
	if !found {
		return ""
	}
	m := p.Municipalities[index]
	return m.Label
}

// NewBuilder creates a builder
func NewBuilder() *MapBuilder {
	return &MapBuilder{
		index:   0,
		indexes: make(map[string]MapIndex, 0),
		target: PostalCodeMap{
			Postalcodes:    make(map[string][]MapIndex, 0),
			Municipalities: make(map[MapIndex]Municipality, 0),
		},
	}
}

// Register create new municipality in the registry and returns its index in the registry to be used in postal code mapping
func (b *MapBuilder) Register(code string, label string) MapIndex {
	index, found := b.indexes[code]
	if found {
		return index
	}
	b.index += 1
	b.indexes[code] = b.index
	b.target.Municipalities[b.index] = Municipality{Code: code, Label: label}
	return b.index
}

// Has test if a Municipality code already exists
// Returns its map index if exists, or false in the boolean code if not
func (b *MapBuilder) Has(municipalityCode string) (MapIndex, bool) {
	index, found := b.indexes[municipalityCode]
	return index, found
}

// AddForPostal associate a postal code to a municipality index
func (b *MapBuilder) AddForPostal(code string, municipalityIdx MapIndex) {
	pp, found := b.target.Postalcodes[code]
	if !found {
		pp = make([]MapIndex, 0, 5)
	}
	pp = append(pp, municipalityIdx)
	pp = unique(pp)
	b.target.Postalcodes[code] = pp
}

// GetMap returns the built registry mapping structure
func (b *MapBuilder) GetMap() *PostalCodeMap {
	return &b.target
}
