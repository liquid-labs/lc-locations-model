package locations

import (
  "strconv"
  "strings"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type Address struct {
  tableName struct{} `sql:"select:addresses_join_locations"`
  Location
  EntityID  EID      `json:"-" sql:",pk"`
  Label     string   `json:"label"`
}

func NewAddress(
    name             string,
    desc             string,
    ownerID          EID,
    publiclyReadable bool,
    address1         string,
    address2         string,
    city             string,
    state            string,
    zip              string,
    entityID         EID,
    label            string) *Address {
  return &Address{
    Location : *NewLocation(name, desc, EID(ownerID), publiclyReadable,
                address1, address2, city, state, zip),
    EntityID : entityID,
    Label    : label,
  }
}

func (a *Address) GetEntityID() EID { return a.EntityID }

func (a *Address) GetLabel() string { return a.Label }
func (a *Address) SetLabel(l string) { a.Label = l }

func (add *Address) Clone() *Address {
  return &Address{struct{}{}, *add.Location.Clone(), add.EntityID, add.Label}
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
      if changeDescs == nil { changeDescs = make([]string, 0, 1) }
      changeDesc = strings.TrimSuffix(changeDesc, `.`) + ` on address ` + strconv.Itoa(i + 1) + `.`
      changeDescs = append(changeDescs, changeDesc)
    }
  }

  return changeDescs
}

type AddressInsertable struct {
  tableName struct{} `sql:"addresses"`
  *Address
  Idx int
}
