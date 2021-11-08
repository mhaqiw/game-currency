DROP TABLE IF EXISTS currency;
CREATE TABLE IF NOT EXISTS currency(
    id bigserial PRIMARY KEY,
    name varchar(225) NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

DROP TABLE IF EXISTS conversion;
CREATE TABLE IF NOT EXISTS conversion(
    id bigserial PRIMARY KEY,
    from_id bigint NOT NULL ,
    to_id bigint NOT NULL ,
    rate DECIMAL(10,5) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (from_id, to_id)
);

ALTER TABLE conversion  add constraint fk_from_id_conversion foreign key (from_id) REFERENCES currency (id);
ALTER TABLE conversion  add constraint fk_to_id_conversion foreign key (to_id) REFERENCES currency (id);
CREATE INDEX idx_conversion_from  ON conversion(from_id);
CREATE INDEX idx_conversion_to  ON conversion(to_id);

INSERT INTO currency (name) values ('Knut'), ('Sickle'), ('Galleon');

