package locations

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type queriesCallback func (*Address, int) []*orm.Query

func (a *Addresses) gatherQueries(f queriesCallback) []*orm.Query {
  qs := make([]*orm.Query, 0, 3*len([]*Address(*a)))
  for i, address := range []*Address(*a) {
    qs = append(qs, f(address, i + 1)...)
  }
  return qs
}

func (a *Addresses) CreateQueries(db orm.DB) []*orm.Query {
  return a.gatherQueries(func (a *Address, i int) []*orm.Query { return a.createQueries(db, i) })
}

// Addresses, as a group, are currently updated by deleting and clearing

// Addresses are not archived

func (a *Addresses) DeleteQueries(db orm.DB) []*orm.Query {
    return a.gatherQueries(func (a *Address, i int) []*orm.Query { return a.deleteQueries(db) })
}

func (a *Addresses) RetrieveByIDRaw(id EID, db orm.DB) error {
  q := db.Model(a).Where(`"address".entity_id=?`, id)
  if err := q.Select(); err != nil {
    return err
  } else {
    return nil
  }
}

// A single Address is not directly managed. Management is at the Addresses level.

func (a *Address) createQueries(db orm.DB, i int) []*orm.Query {
  return append((&a.Location).CreateQueries(db),
    db.Model(&AddressInsertable{Address: a, Idx: i}).ExcludeColumn(LocationFields...))
}

// only update label
func (a *Address) UpdateQueries(db orm.DB) []*orm.Query {
  qs := (&a.Location).UpdateQueries(db)
  q := db.Model(a).
      Set(`label=?label`).
      Where(`"address".id=?id`).
      Where(`"address".entity_id=?entity_id`)
  q.GetModel().Table().SoftDeleteField = nil

  return append(qs, q)
}

// only delete the address part, not location
func (a *Address) deleteQueries(db orm.DB) []*orm.Query {
  q := db.Model(a).Where(`"address".id=?id`).Where(`"address".entity_id=?entity_id`)
  q.GetModel().Table().SoftDeleteField = nil
  return []*orm.Query { q }
}
