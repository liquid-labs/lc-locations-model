package locations_test

import (
  "os"
  "testing"

  // "github.com/go-pg/go/pg/orm"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/strkit/go/strkit"
  . "github.com/Liquid-Labs/lc-users-model/go/users"

  // the package we're testing
  . "github.com/Liquid-Labs/lc-locations-model/go/locations"
)

func (s *AddressesIntegrationSuite) retrieveAddresses(id EID) (*Addresses, error) {
  as := make(Addresses, 0)
  q := rdb.Connect().Model(&as).Where(`"address".entity_id=?`, id)
  _, err := RunRetrieveOp(q, RetrieveOp)
  if err != nil { return nil, err } else { return &as, nil }
}

type AddressesIntegrationSuite struct {
  suite.Suite
  IM *ItemManager
  U  *User
}
func (s *AddressesIntegrationSuite) SetupSuite() {
  s.IM = NewItemManager(rdb.Connect())
  s.IM.AllowUnsafeStateChange = true
  s.U = NewUser(`users`, `test user`, `foo`, strkit.RandString(strkit.LettersAndNumbers, 16), `444-45-3232`, `SSN`, true)
  require.NoError(s.T(), s.IM.CreateRaw(s.U))
}
func TestAddressesesIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(AddressesIntegrationSuite))
  }
}

func (s *AddressesIntegrationSuite) TestAddressesCreateNoOwner() {
  a := NewAddress(`someplace`, `desc`, EID(``), true, `100 Main Str`, `#B`, `City`, `TX`, `78388`, s.U.GetID(), `home`)
  as := &Addresses{a}
  require.NoError(s.T(), s.IM.CreateRaw(as))

  assert.Equal(s.T(), `someplace`, a.GetName())
  assert.Equal(s.T(), `desc`, a.GetDescription())
  assert.Equal(s.T(), a.GetID(), a.GetOwnerID())
  assert.True(s.T(), a.IsPubliclyReadable())
  assert.Equal(s.T(), `100 Main Str`, a.GetAddress1())
  assert.Equal(s.T(), `#B`, a.GetAddress2())
  assert.Equal(s.T(), `City`, a.GetCity())
  assert.Equal(s.T(), `TX`, a.GetState())
  assert.Equal(s.T(), `78388`, a.GetZip())
  assert.Equal(s.T(), s.U.GetID(), a.GetEntityID())
  assert.Equal(s.T(), `home`, a.GetLabel())
}

func (s *AddressesIntegrationSuite) TestAddressUpdate() {
  a := NewAddress(`someplace`, `desc`, EID(``), true, `100 Main Str`, `#B`, `City`, `TX`, `78388`, s.U.GetID(), `home`)
  as := &Addresses{a}
  require.NoError(s.T(), s.IM.CreateRaw(as))

  a.SetName(`Foo`)
  a.SetAddress1(`101 Main Str`)
  a.SetAddress2(`#C`)
  a.SetCity(`Foo`)
  a.SetState(`MN`)
  a.SetZip(`83388`)
  a.SetLabel(`blah`)
  require.NoError(s.T(), s.IM.UpdateRaw(a))

  assert.Equal(s.T(), `Foo`, a.GetName())
  assert.Equal(s.T(), `101 Main Str`, a.GetAddress1())
  assert.Equal(s.T(), `#C`, a.GetAddress2())
  assert.Equal(s.T(), `Foo`, a.GetCity())
  assert.Equal(s.T(), `MN`, a.GetState())
  assert.Equal(s.T(), `83388`, a.GetZip())
  assert.Equal(s.T(), `blah`, a.GetLabel())

  asCopy, err := s.retrieveAddresses(a.GetEntityID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), as, asCopy)
}

func (s *AddressesIntegrationSuite) TestAddressDelete() {
  a := NewAddress(`someplace`, `desc`, EID(``), true, `100 Main Str`, `#B`, `City`, `TX`, `78388`, s.U.GetID(), `home`)
  as := &Addresses{a}
  require.NoError(s.T(), s.IM.CreateRaw(as))
  var e1, l1, a1, e2, l2, a2 int
  rdb.Connect().Query(&e1, "SELECT COUNT(*) FROM entities")
  rdb.Connect().Query(&l1, "SELECT COUNT(*) FROM locations")
  rdb.Connect().Query(&a1, "SELECT COUNT(*) FROM addresses")
  require.NoError(s.T(), s.IM.DeleteRaw(as))
  rdb.Connect().Query(&e2, "SELECT COUNT(*) FROM entities")
  rdb.Connect().Query(&l2, "SELECT COUNT(*) FROM locations")
  rdb.Connect().Query(&a2, "SELECT COUNT(*) FROM addresses")
  assert.Equal(s.T(), e1 - 1, e2)
  assert.Equal(s.T(), l1 - 1, l2)
  assert.Equal(s.T(), a1 - 1, a2)
}
