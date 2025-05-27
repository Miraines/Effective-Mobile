CREATE TABLE people (
                        id           SERIAL PRIMARY KEY,
                        name         TEXT        NOT NULL,
                        surname      TEXT        NOT NULL,
                        patronymic   TEXT,
                        age          INT,
                        gender       TEXT,
                        country_id   CHAR(2),
                        probability  NUMERIC(4,3),
                        created_at   TIMESTAMPTZ DEFAULT now(),
                        updated_at   TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_people_fullname ON people (surname, name, patronymic);
