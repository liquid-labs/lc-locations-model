package locations

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

var LocationFields = append(EntityFields,
  `address1`,
  `address2`,
  `city`,
  `state`,
  `zip`,
  `lat`,
  `lng`)

func (l *Location) CreateQueries(db orm.DB) []*orm.Query {
  return append(
    (&l.Entity).CreateQueries(db),
    db.Model(l).ExcludeColumn(EntityFields...))
}

func (l *Location) UpdateQueries(db orm.DB) []*orm.Query {
  qs := (&l.Entity).UpdateQueries(db)
  q :=  db.Model(l).
    ExcludeColumn(EntityFields...).
    Where(`"location".id=?id`)
  q.GetModel().Table().SoftDeleteField = nil

  return append(qs, q)
}

func (l *Location) ArchiveQueries(db orm.DB) []*orm.Query {
  return (&l.Entity).ArchiveQueries(db)
}

func (l *Location) DeleteQueries(db orm.DB) []*orm.Query {
  q := db.Model(l).
    ExcludeColumn(EntityFields...).
    Where(`"location".id=?id`)
  q.GetModel().Table().SoftDeleteField = nil
  qs := []*orm.Query{q}

  return append(qs, (&l.Entity).DeleteQueries(db)...)
}
