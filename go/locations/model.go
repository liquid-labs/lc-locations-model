package locations

/**
 * Defines basic struct (Go and JSON) and simple data maniplations for locations
 * resource.
 */

import (
  "bytes"
  "database/sql"
  "errors"
  "strconv"
  "strings"

  "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type Location struct {
  LocationID entities.InternalID `json:"locationId"`
  Address1   string              `json:"address1"`
  Address2   string              `json:"address2"`
  City       string              `json:"city"`
  State      string              `json:"state"`
  Zip        string              `json:"zip,string"`
  Lat        sql.NullFloat64     `json:"lat",pg:",use_zero"`
  Lng        sql.NullFloat64     `json:"lng",pg:",use_zero"`
  ChangeDesc []string            `json:"changeDesc,omitempty",sql:"-"`
}

func (loc *Location) Clone() *Location {
  newChangeDesc := make([]string, len(loc.ChangeDesc))
  copy(newChangeDesc, loc.ChangeDesc)
  return &Location{loc.LocationID, loc.Address1, loc.Address2, loc.City, loc.State, loc.Zip, loc.Lat, loc.Lng, newChangeDesc}
}

func (loc *Location) requiredAddressComponents() ([]string) {
  return []string{
    loc.Address1,
    loc.City,
    loc.State,
    loc.Zip,
  }
}

func (loc *Location) latLngComponents() ([]sql.NullFloat64) {
  return []sql.NullFloat64{
    loc.Lat,
    loc.Lng,
  }
}

func (loc *Location) AddressString() (string, error) {
  if !loc.IsAddressComplete() {
    return "", errors.New("Cannot generate address string from incomplete address components.")
  }
  var buffer bytes.Buffer

  buffer.WriteString(loc.Address1)
  if loc.Address2 != `` {
    buffer.WriteString("; ")
    buffer.WriteString(loc.Address2)
  }
  buffer.WriteString(", ")
  buffer.WriteString(loc.City)
  buffer.WriteString(", ")
  buffer.WriteString(loc.State)
  buffer.WriteString(" ")
  buffer.WriteString(loc.Zip)

  return buffer.String(), nil
}

func (loc *Location) IsComplete() (bool) {
  return loc.IsAddressComplete() && loc.IsLatLngComplete()
}

func (loc *Location) IsAddressComplete() (bool) {
  for _, component := range loc.requiredAddressComponents() {
    if component == `` {
      return false
    }
  }
  return true
}

func (loc *Location) IsAddressEmpty() (bool) {
  for _, component := range loc.requiredAddressComponents() {
    if component != `` {
      return false
    }
  }
  return true
}

func (loc *Location) IsLatLngComplete() (bool) {
  for _, component := range loc.latLngComponents() {
    if !component.Valid {
      return false
    }
  }
  return true
}


func (loc *Location) IsLatLngEmpty() (bool) {
  for _, component := range loc.latLngComponents() {
    if component.Valid {
      return false
    }
  }
  return true
}

type Address struct {
  Location
  IDX      int64  `json:"idx"`
  Label    string `json:"label"`
}

func (add *Address) Clone() *Address {
  return &Address{*add.Location.Clone(), add.IDX, add.Label}
}

type Addresses []*Address

func (adds *Addresses) Clone() *Addresses {
  clone := make(Addresses, len(*adds))
  for i, add := range *adds {
    clone[i] = add.Clone()
  }

  return &clone
}

func (adds *Addresses) PromoteChanges(changeDescs []string) ([]string) {
  for i, address := range *adds {
    for _, changeDesc := range address.ChangeDesc {
      changeDesc = strings.TrimSuffix(changeDesc, `.`) + ` on address ` + strconv.Itoa(i + 1) + `.`
      changeDescs = append(changeDescs, changeDesc)
    }
  }

  return changeDescs
}
