-- Active: 1685821706808@@46.36.41.4@5432@sus

/* DROP TABLE IF EXISTS dokoncene;

DROP TABLE IF EXISTS cviceni;

DROP TABLE IF EXISTS slovnik;

DROP TABLE IF EXISTS lekce; */

/* DROP TABLE IF EXISTS uzivatel; */

CREATE TABLE
    IF NOT EXISTS uzivatel (
        id SERIAL PRIMARY KEY,
        jmeno VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(50) NOT NULL UNIQUE,
        heslo VARCHAR(255) NOT NULL,
        klavesnice VARCHAR(10) DEFAULT 'qwertz'
    );

CREATE TABLE
    IF NOT EXISTS overeni (
        jmeno VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(50) NOT NULL UNIQUE,
        heslo VARCHAR(255) NOT NULL,
        kod VARCHAR(5) NOT NULL,
        cas INT
    );

CREATE TABLE
    IF NOT EXISTS zmena_hesla (
        email VARCHAR(50) NOT NULL UNIQUE,
        kod VARCHAR(5) NOT NULL,
        cas INT
    );

CREATE TABLE
    IF NOT EXISTS lekce (
        id SERIAL PRIMARY KEY,
        pismena VARCHAR(25),
        skupina INT,
        klavesnice VARCHAR(10) DEFAULT 'oboje'
    );

CREATE TABLE
    IF NOT EXISTS cviceni (
        id SERIAL PRIMARY KEY,
        typ VARCHAR(20) DEFAULT 'nova',
        lekce_id INT,
        FOREIGN KEY (lekce_id) REFERENCES lekce(id)
    );

CREATE TABLE
    IF NOT EXISTS dokoncene (
        id SERIAL PRIMARY KEY,
        uziv_id INT,
        cviceni_id INT,
        cpm DECIMAL,
        preklepy INT,
        delkatextu INT,
        cas DECIMAL,
        datum TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (uziv_id) REFERENCES uzivatel(id) ON DELETE CASCADE,
        FOREIGN KEY (cviceni_id) REFERENCES cviceni(id),
        CONSTRAINT unikatni UNIQUE(uziv_id, cviceni_id)
    );

CREATE TABLE
    IF NOT EXISTS slovnik (
        id SERIAL PRIMARY KEY,
        slovo VARCHAR(50),
        lekceQWERTZ_id INT,
        lekceQWERTY_id INT
    );

CREATE TABLE
    IF NOT EXISTS texty (
        id SERIAL PRIMARY KEY,
        jmeno VARCHAR(50),
        text1 TEXT,
        text2 TEXT,
        text3 TEXT,
        text4 TEXT,
        text5 TEXT
    );

INSERT INTO lekce (pismena, skupina) VALUES ('fjgh', 1), ('dk', 1), ('sl', 1), ('aů', 1), ('tz', 2), ('ty', 2), ('ru', 2), ('ei', 2), ('wo', 2), ('qpú', 2), ('vb', 3), ('cn', 3), ('yxm', 3), ('zxm', 3), ('žý', 4), ('řá', 4), ('čí', 4), ('ěšé', 4), (',.', 5), ('!?', 5), ('+=-/', 5);
INSERT INTO lekce (pismena, skupina) VALUES ('fjgh', 1), ('dk', 1), ('sl', 1), ('aů', 1), ('tz', 2), ('ty', 2), ('ru', 2), ('ei', 2), ('wo', 2), ('qpú', 2), ('vb', 3), ('cn', 3), ('yxm', 3), ('zxm', 3), ('žý', 4), ('řá', 4), ('čí', 4), ('ěšé', 4), ('Zbylá diakritika', 4), ('Shift', 5);
