DROP TABLE IF EXISTS dokoncene;
DROP TABLE IF EXISTS cviceni;
DROP TABLE IF EXISTS lekce;
DROP TABLE IF EXISTS uzivatel;

CREATE TABLE IF NOT EXISTS uzivatel (
    id SERIAL PRIMARY KEY, 
    jmeno VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    heslo VARCHAR(255) NOT NULL
);

CREATE TABLE lekce (
    id SERIAL PRIMARY KEY,
    pismena VARCHAR(10),
    skupina INT
);

CREATE TABLE cviceni (
    id SERIAL PRIMARY KEY,
    typ VARCHAR(20) DEFAULT 'nova',
    lekce_id INT,
    FOREIGN KEY (lekce_id)
        REFERENCES lekce(id)
);

CREATE TABLE dokoncene (
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

INSERT INTO lekce (pismena, skupina) VALUES ('fjgh', 1), ('dk', 1), ('sl', 1), ('aů', 1), ('tz', 2);
INSERT INTO cviceni (lekce_id) VALUES (1), (1), (1), (2), (2);