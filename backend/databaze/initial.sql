-- Active: 1685821706808@@46.36.41.4@5432@pavouk

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
        pismena VARCHAR(10),
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
INSERT INTO cviceni (lekce_id, typ) 
VALUES (1, 'nova'), (1, 'nova'), (1, 'nova'), (1, 'nova'),
(2, 'nova'), (2, 'nova'), (2, 'nova'), (2, 'naucena'), (2, 'naucena'),
(3, 'nova'), (3, 'nova'), (3, 'nova'), (3, 'naucena'), (3, 'naucena'), 
(4, 'nova'), (4, 'nova'), (4, 'nova'), (4, 'naucena'), (4, 'naucena'), (4, 'slova'),
(5, 'nova'), (5, 'nova'), (5, 'naucena'), (5, 'slova'), (5, 'slova'), 
(6, 'nova'), (6, 'nova'), (6, 'naucena'), (6, 'slova'), (6, 'slova'), 
(7, 'nova'), (7, 'nova'), (7, 'naucena'), (7, 'slova'), (7, 'slova'),
(8, 'nova'), (8, 'nova'), (8, 'naucena'), (8, 'slova'), (8, 'slova'),
(9, 'nova'), (9, 'nova'), (9, 'naucena'), (9, 'naucena'), (9, 'slova'), (9, 'slova'),
(10, 'nova'), (10, 'nova'), (10, 'naucena'), (10, 'slova'), (10, 'slova'),
(11, 'nova'), (11, 'nova'), (11, 'naucena'), (11, 'slova'), (11, 'slova'),
(12, 'nova'), (12, 'nova'), (12, 'naucena'),  (12, 'naucena'), (12, 'slova'), (12, 'slova'),
(13, 'nova'), (13, 'nova'), (13, 'naucena'), (13, 'slova'), (13, 'slova'),
(14, 'nova'), (14, 'nova'), (14, 'naucena'), (14, 'slova'), (14, 'slova'),
(15, 'nova'), (15, 'nova'), (15, 'naucena'), (15, 'slova'), (15, 'slova'),
(16, 'nova'), (16, 'nova'), (16, 'naucena'), (16, 'naucena'), (16, 'slova'), (16, 'slova'),
(17, 'nova'), (17, 'nova'), (17, 'naucena'),
(18, 'nova'), (18, 'nova'), (18, 'naucena'),
(19, 'nova'), (19, 'nova'), (19, 'naucena');
/* SELECT * FROM slovnik ORDER BY lekce_id; */