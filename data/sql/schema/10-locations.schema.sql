CREATE TABLE locations (
  id       UUID,
  address1 VARCHAR(255) NOT NULL,
  address2 VARCHAR(255),
  city     VARCHAR(255) NOT NULL,
  state    VARCHAR(2) NOT NULL,
  zip      VARCHAR(12) NOT NULL,
  lat      DECIMAL(9,7),
  lng      DECIMAL(10,7),

  CONSTRAINT locations_key PRIMARY KEY ( id ),
  CONSTRAINT locations_refs_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);

CREATE VIEW locations_join_entities AS
  SELECT e.*,
    l.address1, l.address2, l.city, l.state, l.zip, l.lat, l.Lng
  FROM locations l
    JOIN entities e ON l.id=e.id;

CREATE TABLE addresses (
  id        UUID NOT NULL,
  entity_id UUID NOT NULL,
  label     VARCHAR(64),
  idx       INT NOT NULL CHECK(idx >= 0 AND idx < 100),

  CONSTRAINT addresses_key PRIMARY KEY (id, entity_id),
  CONSTRAINT addresses_refs_locations FOREIGN KEY ( id ) REFERENCES locations ( id ),
  CONSTRAINT addresses_refs_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);

CREATE INDEX addresses_by_location ON addresses ( entity_id );

CREATE VIEW addresses_join_locations AS
  SELECT e.*,
      l.address1, l.address2, l.city, l.state, l.zip, l.lat, l.lng,
      a.entity_id, a.label, a.idx
    FROM addresses a
      JOIN locations l ON l.id=a.id
      JOIN entities e ON l.id=e.id;
