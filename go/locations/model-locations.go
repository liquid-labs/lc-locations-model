package locations

/**
 * Defines basic struct (Go and JSON) and simple data maniplations for locations
 * resource.
 */

import (
  "bytes"
  "database/sql"
  "errors"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

const LocationsResName = ResourceName(`locations`)

type Location struct {
  Entity
  Address1   string          `json:"address1"`
  Address2   string          `json:"address2"`
  City       string          `json:"city"`
  State      string          `json:"state"`
  Zip        string          `json:"zip"`
  Lat        sql.NullFloat64 `json:"lat" pg:",use_zero"`
  Lng        sql.NullFloat64 `json:"lng" pg:",use_zero"`
  ChangeDesc []string        `json:"changeDesc,omitempty" sql:"-"`
}

func NewLocation(
    name             string,
    desc             string,
    ownerID          EID,
    publiclyReadable bool,
    address1         string,
    address2         string,
    city             string,
    state            string,
    zip              string) *Location {
  return &Location{
    *NewEntity(LocationsResName, name, desc, ownerID, publiclyReadable),
    address1, address2, city, state, zip,
    sql.NullFloat64{Valid:false}, sql.NullFloat64{Valid:false}, nil}
}

func (l *Location) GetAddress1() string { return l.Address1 }
func (l *Location) SetAddress1(a1 string) { l.Address1 = a1 }

func (l *Location) GetAddress2() string { return l.Address2 }
func (l *Location) SetAddress2(a2 string) { l.Address2 = a2 }

func (l *Location) GetCity() string { return l.City }
func (l *Location) SetCity(c string) { l.City = c }

func (l *Location) GetState() string { return l.State }
func (l *Location) SetState(s string) { l.State = s }

func (l *Location) GetZip() string { return l.Zip }
func (l *Location) SetZip(z string) { l.Zip = z }

func (l *Location) GetLat() sql.NullFloat64 { return l.Lat }

func (l *Location) GetLng() sql.NullFloat64 { return l.Lng }

func (loc *Location) Clone() *Location {
  newChangeDesc := ([]string)(nil)
  if loc.ChangeDesc != nil {
    newChangeDesc := make([]string, len(loc.ChangeDesc))
    copy(newChangeDesc, loc.ChangeDesc)
  }
  return &Location{*(&loc.Entity).Clone(), loc.Address1, loc.Address2, loc.City, loc.State, loc.Zip, loc.Lat, loc.Lng, newChangeDesc}
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
