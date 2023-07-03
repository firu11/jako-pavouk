DROP TABLE IF EXISTS dokoncene;
DROP TABLE IF EXISTS cviceni;
DROP TABLE IF EXISTS slovnik;
DROP TABLE IF EXISTS lekce;
DROP TABLE IF EXISTS uzivatel;

CREATE TABLE IF NOT EXISTS uzivatel (
    id SERIAL PRIMARY KEY, 
    jmeno VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    heslo VARCHAR(255) NOT NULL,
    daystreak INT DEFAULT 0,
    posledniden DATE DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS lekce (
    id SERIAL PRIMARY KEY,
    pismena VARCHAR(10),
    skupina INT
);

CREATE TABLE IF NOT EXISTS cviceni (
    id SERIAL PRIMARY KEY,
    typ VARCHAR(20) DEFAULT 'nova',
    lekce_id INT,
    FOREIGN KEY (lekce_id)
        REFERENCES lekce(id)
);

CREATE TABLE IF NOT EXISTS dokoncene (
    id SERIAL PRIMARY KEY,
    uziv_id INT,
    cviceni_id INT,
    cpm DECIMAL,
    preklepy INT,
    FOREIGN KEY (uziv_id)
        REFERENCES uzivatel(id)
        ON DELETE CASCADE,
    FOREIGN KEY (cviceni_id)
        REFERENCES cviceni(id),
    CONSTRAINT unikatni UNIQUE(uziv_id, cviceni_id)
);

CREATE TABLE IF NOT EXISTS slovnik (
    id SERIAL PRIMARY KEY,
    slovo VARCHAR(50),
    lekce_id INT
);

INSERT INTO lekce (pismena, skupina) VALUES ('fjgh', 1), ('dk', 1), ('sl', 1), ('aů', 1), ('tz', 2), ('ru', 2), ('ei', 2), ('wo', 2), ('qpú', 2), ('vb', 3), ('cn', 3), ('yxm', 3), ('žý', 4), ('řá', 4), ('čí', 4), ('ěšé', 4), (',.', 5), ('!?', 5), ('+=-/', 5);
INSERT INTO cviceni (lekce_id, typ) VALUES (1, 'nova'), (1, 'nova'), (1, 'nova'), (1, 'nova'), (2, 'nova'), (2, 'nova'), (2, 'naucena'), (3, 'nova'), (3, 'nova'), (3, 'naucena'), (4, 'nova'), (4, 'naucena'), (4, 'slova'), (5, 'nova'), (5, 'slova');

/* SELECT * FROM slovnik ORDER BY lekce_id; */
