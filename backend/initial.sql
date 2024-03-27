--
-- PostgreSQL database dump
--

-- Dumped from database version 14.11 (Ubuntu 14.11-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.11 (Ubuntu 14.11-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cviceni; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cviceni (
    id integer NOT NULL,
    typ character varying(20) DEFAULT 'nova'::character varying,
    lekce_id integer
);


--
-- Name: cviceni_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.cviceni_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cviceni_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.cviceni_id_seq OWNED BY public.cviceni.id;


--
-- Name: dokoncene; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.dokoncene (
    id integer NOT NULL,
    uziv_id integer,
    cviceni_id integer,
    cpm numeric,
    preklepy integer,
    cas numeric,
    den timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    delka_textu integer DEFAULT 50
);


--
-- Name: dokoncene_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.dokoncene_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: dokoncene_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.dokoncene_id_seq OWNED BY public.dokoncene.id;


--
-- Name: lekce; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lekce (
    id integer NOT NULL,
    pismena character varying(25),
    skupina integer,
    klavesnice character varying(10) DEFAULT 'oboje'::character varying
);


--
-- Name: lekce_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.lekce_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: lekce_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.lekce_id_seq OWNED BY public.lekce.id;


--
-- Name: navstevnost; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.navstevnost (
    id integer NOT NULL,
    den date,
    pocet integer
);


--
-- Name: navstevnost_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.navstevnost_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: navstevnost_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.navstevnost_id_seq OWNED BY public.navstevnost.id;


--
-- Name: overeni; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.overeni (
    jmeno character varying(50) NOT NULL,
    email character varying(50) NOT NULL,
    heslo character varying(255) NOT NULL,
    kod character varying(5) NOT NULL,
    cas integer
);


--
-- Name: slovnik; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.slovnik (
    id integer NOT NULL,
    slovo character varying(50),
    lekceqwertz_id integer,
    lekceqwerty_id integer
);


--
-- Name: slovnik_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.slovnik_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: slovnik_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.slovnik_id_seq OWNED BY public.slovnik.id;


--
-- Name: slovnik_programator; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.slovnik_programator (
    id integer NOT NULL,
    slovo character varying(50)
);


--
-- Name: slovnik_programator_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.slovnik_programator_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: slovnik_programator_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.slovnik_programator_id_seq OWNED BY public.slovnik_programator.id;


--
-- Name: texty; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.texty (
    id integer NOT NULL,
    jmeno character varying(50),
    text1 text,
    text2 text,
    text3 text,
    text4 text,
    text5 text,
    text6 text,
    text7 text,
    text8 text,
    text9 text,
    text10 text
);


--
-- Name: texty_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.texty_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: texty_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.texty_id_seq OWNED BY public.texty.id;


--
-- Name: uzivatel; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uzivatel (
    id integer NOT NULL,
    jmeno character varying(50) NOT NULL,
    email character varying(50) NOT NULL,
    heslo character varying(255) NOT NULL,
    klavesnice character varying(10) DEFAULT 'qwertz'::character varying
);


--
-- Name: uzivatel_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.uzivatel_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uzivatel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.uzivatel_id_seq OWNED BY public.uzivatel.id;


--
-- Name: vety; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vety (
    id integer NOT NULL,
    veta text NOT NULL,
    delka integer NOT NULL
);


--
-- Name: vety_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.vety_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: vety_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.vety_id_seq OWNED BY public.vety.id;


--
-- Name: zmena_hesla; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.zmena_hesla (
    email character varying(50) NOT NULL,
    kod character varying(5) NOT NULL,
    cas integer
);


--
-- Name: cviceni id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cviceni ALTER COLUMN id SET DEFAULT nextval('public.cviceni_id_seq'::regclass);


--
-- Name: dokoncene id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokoncene ALTER COLUMN id SET DEFAULT nextval('public.dokoncene_id_seq'::regclass);


--
-- Name: lekce id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lekce ALTER COLUMN id SET DEFAULT nextval('public.lekce_id_seq'::regclass);


--
-- Name: navstevnost id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.navstevnost ALTER COLUMN id SET DEFAULT nextval('public.navstevnost_id_seq'::regclass);


--
-- Name: slovnik id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.slovnik ALTER COLUMN id SET DEFAULT nextval('public.slovnik_id_seq'::regclass);


--
-- Name: slovnik_programator id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.slovnik_programator ALTER COLUMN id SET DEFAULT nextval('public.slovnik_programator_id_seq'::regclass);


--
-- Name: texty id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.texty ALTER COLUMN id SET DEFAULT nextval('public.texty_id_seq'::regclass);


--
-- Name: uzivatel id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uzivatel ALTER COLUMN id SET DEFAULT nextval('public.uzivatel_id_seq'::regclass);


--
-- Name: vety id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vety ALTER COLUMN id SET DEFAULT nextval('public.vety_id_seq'::regclass);


--
-- Name: cviceni cviceni_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cviceni
    ADD CONSTRAINT cviceni_pkey PRIMARY KEY (id);


--
-- Name: dokoncene dokoncene_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokoncene
    ADD CONSTRAINT dokoncene_pkey PRIMARY KEY (id);


--
-- Name: lekce lekce_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lekce
    ADD CONSTRAINT lekce_pkey PRIMARY KEY (id);


--
-- Name: navstevnost navstevnost_den_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.navstevnost
    ADD CONSTRAINT navstevnost_den_key UNIQUE (den);


--
-- Name: navstevnost navstevnost_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.navstevnost
    ADD CONSTRAINT navstevnost_pkey PRIMARY KEY (id);


--
-- Name: overeni overeni_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.overeni
    ADD CONSTRAINT overeni_email_key UNIQUE (email);


--
-- Name: overeni overeni_jmeno_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.overeni
    ADD CONSTRAINT overeni_jmeno_key UNIQUE (jmeno);


--
-- Name: slovnik slovnik_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.slovnik
    ADD CONSTRAINT slovnik_pkey PRIMARY KEY (id);


--
-- Name: slovnik_programator slovnik_programator_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.slovnik_programator
    ADD CONSTRAINT slovnik_programator_pkey PRIMARY KEY (id);


--
-- Name: texty texty_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.texty
    ADD CONSTRAINT texty_pkey PRIMARY KEY (id);


--
-- Name: dokoncene unikatni; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokoncene
    ADD CONSTRAINT unikatni UNIQUE (uziv_id, cviceni_id);


--
-- Name: uzivatel uzivatel_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uzivatel
    ADD CONSTRAINT uzivatel_email_key UNIQUE (email);


--
-- Name: uzivatel uzivatel_jmeno_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uzivatel
    ADD CONSTRAINT uzivatel_jmeno_key UNIQUE (jmeno);


--
-- Name: uzivatel uzivatel_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uzivatel
    ADD CONSTRAINT uzivatel_pkey PRIMARY KEY (id);


--
-- Name: vety vety_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vety
    ADD CONSTRAINT vety_pkey PRIMARY KEY (id);


--
-- Name: zmena_hesla zmena_hesla_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.zmena_hesla
    ADD CONSTRAINT zmena_hesla_email_key UNIQUE (email);


--
-- Name: cviceni cviceni_lekce_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cviceni
    ADD CONSTRAINT cviceni_lekce_id_fkey FOREIGN KEY (lekce_id) REFERENCES public.lekce(id);


--
-- Name: dokoncene dokoncene_cviceni_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokoncene
    ADD CONSTRAINT dokoncene_cviceni_id_fkey FOREIGN KEY (cviceni_id) REFERENCES public.cviceni(id);


--
-- Name: dokoncene dokoncene_uziv_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokoncene
    ADD CONSTRAINT dokoncene_uziv_id_fkey FOREIGN KEY (uziv_id) REFERENCES public.uzivatel(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

COPY public.lekce (id, pismena, skupina, klavesnice) FROM stdin;
1	fjgh	1	oboje
2	dk	1	oboje
3	sl	1	oboje
4	aů	1	oboje
7	ru	2	oboje
8	ei	2	oboje
9	wo	2	oboje
10	qpú	2	oboje
11	vb	3	oboje
12	cn	3	oboje
15	žý	4	oboje
16	řá	4	oboje
17	čí	4	oboje
18	ěšé	4	oboje
5	tz	2	qwertz
6	ty	2	qwerty
13	yxm	3	qwertz
14	zxm	3	qwerty
19	zbylá diakritika	4	oboje
20	velká písmena (shift)	5	oboje
21	čísla	5	oboje
22	závorky	6	oboje
23	operátory	6	oboje
\.


COPY public.cviceni (id, typ, lekce_id) FROM stdin;
1	nova	13
2	nova	13
3	nova	13
4	naucena	13
5	slova	13
6	slova	13
7	nova	17
8	nova	17
9	naucena	17
10	slova	17
11	slova	17
12	nova	18
13	nova	18
14	nova	18
15	naucena	18
16	slova	18
17	slova	18
18	nova	2
19	nova	2
20	nova	2
21	naucena	2
22	naucena	2
23	nova	4
24	nova	4
25	nova	4
26	naucena	4
27	naucena	4
28	slova	4
29	nova	6
30	nova	6
31	naucena	6
32	slova	6
33	slova	6
34	nova	11
35	nova	11
36	naucena	11
37	slova	11
38	slova	11
39	nova	12
40	nova	12
41	naucena	12
42	slova	12
43	slova	12
50	nova	3
51	nova	3
52	nova	3
53	naucena	3
54	naucena	3
55	nova	7
56	nova	7
57	naucena	7
58	slova	7
59	slova	7
60	nova	8
61	nova	8
62	naucena	8
63	slova	8
64	slova	8
65	nova	16
66	nova	16
67	naucena	16
68	slova	16
69	slova	16
73	nova	15
74	nova	15
75	naucena	15
76	slova	15
77	slova	15
78	nova	1
79	nova	1
80	nova	1
81	nova	1
82	nova	5
83	nova	5
84	naucena	5
85	slova	5
86	slova	5
87	nova	9
88	nova	9
89	naucena	9
90	slova	9
91	slova	9
92	nova	10
93	nova	10
94	nova	10
95	naucena	10
96	slova	10
97	slova	10
98	nova	14
99	nova	14
100	nova	14
101	naucena	14
102	slova	14
103	slova	14
104	nova	19
105	naucena	19
106	slova	19
107	slova	19
108	naucena	20
109	naucena	20
110	slova	20
111	slova	20
120	naucena	21
121	slova	21
122	slova	21
123	programator	22
124	programator	22
125	programator	22
126	programator	23
127	programator	23
128	programator	23
\.


COPY public.vety (id, veta, delka) FROM stdin;
1	Když princeznu spatřila, neváhala a na spánek se jí zeptala.	60
2	Přitom pořád dávala pozor na svého malého svěřence v tričku.	60
3	A ačkoliv byl o pár roků starší jak Péťa, rozuměli si spolu.	60
4	Pak tatínka myslivce napadlo, že nejvíc toho uvidí z posedu.	60
5	Rozpoutala se sněhová bouře a chatrný králičí kotec zničila.	60
6	A taky vyrábění pomohlo k tomu, aby zachránil ostatní ptáky.	60
7	Zahrady, louky, parky, ale i dětská hřiště jimi byla posetá.	60
8	Večer se však neúprosně blížil a s ním i jarní noční mrazík.	60
9	Všechny věci patřily Jance, jejímu pejskovi a jejím rodičům.	60
10	Když tu najednou nedaleko v trávě slyšela tiché pofňukávání.	60
11	Proto polednice ztratila mezi těmi veselými lidmi svojí moc.	60
12	A tak zkusil prozkoumat hned první chatu, na kterou narazil.	60
13	Rochnila si tam spokojeně už dlouhý čas a nic jim nechybělo.	60
14	Když ji zbyl čas, jezdila se svým koněm Ucháčem na vyjížďky.	60
15	Mohli by proto pohádkové bytosti z lesa vyhnat nebo uvěznit.	60
16	A někdy někde se objevil i soplíkový rampouch někomu u nosu.	60
17	Přesto ví o všem, co se kolem děje, a vždycky někomu pomůže.	60
18	Z ničehož nic spatřil obrovský blesk a poté i zaslechl ránu.	60
19	Seskočila z parapetu a šla za pejskem tam, kde myš měla být.	60
20	Ráno přišel král, přepočítal jablíčka a zase jedno scházelo.	60
21	Když jí strýček prováděl svou farmou, byla naprosto unešená.	60
22	Často jezdil do cizích zemí, aby byl na světě pořádek a mír.	60
23	A právě k těm koledám se váže symbolika velikonočních barev.	60
24	Také si uvědomila, že už je večer a ona měla být dávno doma.	60
25	Už se těšila, jak si to pejsek přiběhne sníst, ale on nikde.	60
26	To je něco jako volná příroda, ale pořád na ní dohlíží lidé.	60
27	Sem tam do studánky spadl lístek ze stromu, ale to nevadilo.	60
28	Nemohla vstát, nemohla ani volat o pomoc, protože byla němá.	60
29	Dávala si přitom velký pozor, aby se neopařila horkou vodou.	60
30	Byla sice krásná, ale oni se její krásy a kouzel velmi báli.	60
31	Čistou a v přepychových šatech ji ani její sestry nepoznaly.	60
32	Po svatebních koláčích se symbolicky podávala hrachová kaše.	60
33	Protože stonek pouštěl zelenou barvu, sloužil jako propiska.	60
34	Byla by pak jako ten zlý pes z pohádky O pejskovi a kočičce.	60
35	V obou rukách pevně držela zrcátko, kterým mířila na slunce.	60
36	Modrého ptáka se podařilo proměnit zpět na prince Alexandra.	60
37	Smutná sedla nedaleko díry a přemýšlela o nenadálé události.	60
38	Spadli tak nešťastně, že koník Blesk spadl rovnou na prince.	60
39	Nakonec Adélku vodníkovi vrátil a on bouřky a déšť zastavil.	60
40	Na Modré pondělí se podle tradic nemá pracovat, ale uklízet.	60
41	Babička jí všechno vysvětlovala a Anička se nestačila divit.	60
42	Vodník došel až ke dveřím, a viděl, že u nich potůček končí.	60
43	Dole na palouku byla vysoká tráva a rosa na ní jako granáty.	60
44	Políbila ho na čelo a Káj zapomněl na svou nejmilejší Gerdu.	60
45	Byl to totiž čaj, který mu jeho žena, královna, vždy vařila.	60
46	Karlík s Lotkou sklopili hlavy k zemi a mamince se omluvili.	60
47	Doma to náramně prokouklo a vypadalo to tam mnohem útulněji.	60
48	Pořád číhal a čekal, až bude slůně co nejdál od své skupiny.	60
49	I kdyby byl poslední, chtěl aspoň uběhnout těch pár koleček.	60
50	Víle Adélce se žádná nemohla vyrovnat, v kráse, ani v tanci.	60
\.

COPY public.slovnik (id, slovo, lekceqwertz_id, lekceqwerty_id) FROM stdin;
1	abnormálně	18	18
2	abnormální	17	17
3	abraham	13	14
4	abrahama	13	14
5	absence	12	12
6	absenci	12	12
7	absencí	17	17
8	absolutně	18	18
9	absolutní	17	17
10	absolutního	17	17
11	absolutních	17	17
12	absolutním	17	17
13	absolvent	12	12
14	absolventa	12	12
15	absolventem	13	14
16	absolventi	12	12
17	absolventů	12	12
18	absolventy	13	12
19	absolvoval	11	11
20	absolvovala	11	11
21	absolvovali	11	11
22	absolvování	17	17
23	absolvovat	11	11
24	absolvuje	11	11
25	absolvují	17	17
26	absorbovat	11	11
27	absorbuje	11	11
28	absorpce	12	12
29	absorpční	17	17
30	abstinence	12	12
31	abstinenci	12	12
32	abstrakce	12	12
33	abstrakci	12	12
34	abstraktní	17	17
35	abstraktního	17	17
36	abstraktních	17	17
37	abstraktním	17	17
38	absurdity	13	11
39	absurdně	18	18
40	absurdní	17	17
41	absurdních	17	17
42	aby	13	11
43	abych	13	12
44	abychom	13	14
45	abys	13	11
46	abysme	13	14
47	abyste	13	11
48	academia	13	14
49	academy	13	14
50	acer	12	12
51	active	12	12
52	áčka	17	17
53	áčko	17	17
54	ačkoli	17	17
55	ačkoliv	17	17
56	adam	13	14
57	adama	13	14
58	adame	13	14
59	adamec	13	14
60	adámek	16	16
61	adamem	13	14
62	adamovi	13	14
63	adaptace	12	12
64	adaptaci	12	12
65	adaptací	17	17
66	adaptér	18	18
67	adaptéru	18	18
68	adaptivní	17	17
69	adaptovat	11	11
70	adekvátně	18	18
71	adekvátní	17	17
72	adéla	18	18
73	adele	8	8
74	adeline	12	12
75	adélka	18	18
76	adeptů	10	10
77	adepty	13	10
78	adhd	4	4
79	adidas	8	8
80	administrace	13	14
81	administraci	13	14
82	administrativa	13	14
83	administrativě	18	18
84	administrativně	18	18
85	administrativní	17	17
86	administrativních	17	17
87	administrativu	13	14
88	administrativy	13	14
89	administrátor	16	16
90	admirál	16	16
91	adolf	9	9
92	adolfa	9	9
93	adopce	12	12
94	adopci	12	12
95	adoptivní	17	17
96	adoptovat	11	11
97	adrenalin	12	12
98	adrenalinové	18	18
99	adrenalinu	12	12
100	adresa	8	8
101	adresář	16	16
102	adresáře	16	16
103	adresáta	16	16
104	adrese	8	8
105	adresou	9	9
106	adresu	8	8
107	adresy	13	8
108	adrian	12	12
109	adriana	12	12
110	advent	12	12
111	adventní	17	17
112	adventu	12	12
113	advokacie	12	12
114	advokát	16	16
115	advokáta	16	16
116	advokátem	16	16
117	advokáti	16	16
118	advokátka	16	16
119	advokátní	17	17
120	advokátů	16	16
121	advokáty	16	16
122	aedeagus	8	8
123	aero	9	9
124	aerobiku	11	11
125	aerobní	17	17
126	aerolinek	12	12
127	aerolinie	12	12
128	aerolinií	17	17
129	aerolinky	13	12
130	afektivní	17	17
131	aféra	18	18
132	aféru	18	18
133	aféry	18	18
134	aféře	18	18
135	afghánistán	16	16
136	afghánistánu	16	16
137	afghánské	18	18
138	afghánských	16	16
139	africe	12	12
140	africká	16	16
141	africké	18	18
142	afrického	18	18
143	africkém	18	18
144	africkou	12	12
145	africký	15	15
146	afrických	15	15
147	africkým	15	15
148	afričané	18	18
149	afrika	8	8
150	afrikou	9	9
151	afriku	8	8
152	afriky	13	8
153	afrodita	9	9
154	agáta	16	16
155	agent	12	12
156	agenta	12	12
157	agentem	13	14
158	agenti	12	12
159	agentovi	12	12
160	agentů	12	12
161	agentura	12	12
162	agenturu	12	12
163	agentury	13	12
164	agentuře	16	16
165	agenty	13	12
166	agrese	8	8
167	agresi	8	8
168	agresivita	11	11
169	agresivitě	18	18
170	agresivitou	11	11
171	agresivitu	11	11
172	agresivity	13	11
173	agresivně	18	18
174	agresivnější	18	18
175	agresivní	17	17
176	agresivního	17	17
177	agresivních	17	17
178	agresivním	17	17
179	aha	4	4
180	ahmed	13	14
181	ahoj	9	9
182	ach	12	12
183	aids	8	8
184	airbagy	13	11
185	airways	13	9
186	akademická	16	16
187	akademické	18	18
188	akademického	18	18
189	akademickém	18	18
190	akademickou	13	14
191	akademický	15	15
192	akademických	15	15
193	akademie	13	14
194	akademii	13	14
195	akademií	17	17
196	akademiků	13	14
197	akce	12	12
198	akcelerace	12	12
199	akceleraci	12	12
200	akcemi	13	14
201	akcent	12	12
202	akcentem	13	14
203	akceptovat	12	12
204	akceptuje	12	12
205	akci	12	12
206	akcí	17	17
207	akcie	12	12
208	akciemi	13	14
209	akcích	17	17
210	akcii	12	12
211	akcií	17	17
212	akcím	17	17
213	akcionář	16	16
214	akcionáře	16	16
215	akcionářem	16	16
216	akcionáři	16	16
217	akcionářů	16	16
218	akcionářům	16	16
219	akciová	16	16
220	akciové	18	18
221	akciovou	12	12
222	akciových	15	15
223	akční	17	17
224	akčního	17	17
225	akčních	17	17
226	akčním	17	17
227	akné	18	18
228	akorát	16	16
229	akord	9	9
230	akordy	13	9
231	akreditace	12	12
232	akreditaci	12	12
233	akropolis	10	10
234	akrů	7	7
235	akt	5	6
236	aktem	13	14
237	aktér	18	18
238	aktéra	18	18
239	aktérem	18	18
240	aktérů	18	18
241	aktéry	18	18
242	aktéři	18	18
243	aktiv	11	11
244	aktiva	11	11
245	aktivace	12	12
246	aktivaci	12	12
247	aktivační	17	17
248	aktivista	11	11
249	aktivisté	18	18
250	aktivistů	11	11
251	aktivisty	13	11
252	aktivit	11	11
253	aktivita	11	11
254	aktivitách	16	16
255	aktivitám	16	16
256	aktivitami	13	14
257	aktivitě	18	18
258	aktivitou	11	11
259	aktivitu	11	11
260	aktivity	13	11
261	aktivně	18	18
262	aktivnější	18	18
263	aktivní	17	17
264	aktivního	17	17
265	aktivních	17	17
266	aktivním	17	17
267	aktivními	17	17
268	aktivnímu	17	17
269	aktivovat	11	11
270	aktivuje	11	11
271	aktivují	17	17
272	aktovku	11	11
273	aktovky	13	11
274	aktualizace	12	14
275	aktualizaci	12	14
276	aktualizací	17	17
277	aktualizované	18	18
278	aktualizovat	11	14
279	aktuálně	18	18
280	aktuální	17	17
281	aktuálního	17	17
282	aktuálních	17	17
283	aktuálním	17	17
284	aktuálními	17	17
285	aktuálnímu	17	17
286	akty	13	6
287	akumulace	13	14
288	akumulaci	13	14
289	akumulací	17	17
290	akumulační	17	17
291	akumulátor	16	16
292	akumulátoru	16	16
293	akumulátorů	16	16
294	akumulátory	16	16
295	akustické	18	18
296	akustického	18	18
297	akustickou	12	12
298	akustický	15	15
299	akustických	15	15
300	akutní	17	17
301	akutního	17	17
302	akutních	17	17
303	akutním	17	17
304	akvária	16	16
305	akváriích	17	17
306	akváriu	16	16
307	akvárium	16	16
308	akvinský	15	15
309	akvizic	12	14
310	akvizice	12	14
311	akvizici	12	14
312	alan	12	12
313	alana	12	12
314	alarm	13	14
315	alarmující	17	17
316	alba	11	11
317	albánie	16	16
318	albánii	16	16
319	albatros	11	11
320	albem	13	14
321	albert	11	11
322	alberta	11	11
323	alberto	11	11
324	albertovi	11	11
325	albrecht	12	12
326	albrechta	12	12
327	albu	11	11
328	album	13	14
329	alderman	13	14
330	aldo	9	9
331	ale	8	8
332	alegorie	9	9
333	alej	8	8
334	aleje	8	8
335	aleji	8	8
336	alejí	17	17
337	alena	12	12
338	aleně	18	18
339	alenka	12	12
340	alenou	12	12
341	alenu	12	12
342	aleny	13	12
343	alergeny	13	12
344	alergické	18	18
345	alergie	8	8
346	alergii	8	8
347	alergií	17	17
348	alergiky	13	8
349	alespoň	19	19
350	alessandro	12	12
351	aleš	18	18
352	aleše	18	18
353	alešem	18	18
354	alex	13	14
355	alexa	13	14
356	alexander	13	14
357	alexandr	13	14
358	alexandra	13	14
359	alexandrem	13	14
360	alexandrie	13	14
361	alexandrii	13	14
362	alexe	13	14
363	alexej	13	14
364	alexeje	13	14
365	alexis	13	14
366	alfa	4	4
367	alfons	12	12
368	alfonse	12	12
369	alfred	8	8
370	alfréd	18	18
371	alfreda	8	8
372	alfréda	18	18
373	algoritmu	13	14
374	algoritmů	13	14
375	algoritmus	13	14
376	algoritmy	13	14
377	alchymie	13	14
378	alchymii	13	14
379	alchymista	13	14
380	alchymisté	18	18
381	ali	8	8
382	aliance	12	12
383	alianci	12	12
384	alias	8	8
385	alibi	11	11
386	alice	12	12
387	alici	12	12
388	alicí	17	17
389	aliho	9	9
390	alimenty	13	14
391	alina	12	12
392	aljašce	18	18
393	aljašky	18	18
394	alkalické	18	18
395	alkalických	15	15
396	alkohol	9	9
397	alkoholem	13	14
398	alkoholické	18	18
399	alkoholických	15	15
400	alkoholik	9	9
401	alkoholiků	9	9
402	alkoholismu	13	14
403	alkoholismus	13	14
404	alkoholu	9	9
405	alláh	16	16
406	alláha	16	16
407	allan	12	12
408	allen	12	12
409	allena	12	12
410	allianz	12	14
411	almužnu	15	15
412	alobalu	11	11
413	aloe	9	9
414	alois	9	9
415	aloise	9	9
416	alokace	12	12
417	alonso	12	12
418	alpách	16	16
419	alpské	18	18
420	alpských	15	15
421	alpy	13	10
422	altán	16	16
423	altánku	16	16
424	altánu	16	16
425	alter	8	8
426	alternativ	12	12
427	alternativa	12	12
428	alternativní	17	17
429	alternativního	17	17
430	alternativních	17	17
431	alternativním	17	17
432	alternativou	12	12
433	alternativu	12	12
434	alternativy	13	12
435	alyson	13	12
436	alzheimerova	13	14
437	alzheimerově	18	18
438	alzheimerovou	13	14
439	alzheimerovy	13	14
440	alžběta	18	18
441	alžbětě	18	18
442	alžbětu	18	18
443	alžběty	18	18
444	amálie	16	16
445	amancio	13	14
446	amanda	13	14
447	amatér	18	18
448	amatérské	18	18
449	amatérského	18	18
450	amatérský	18	18
451	amatérských	18	18
452	amatéry	18	18
453	amatéři	18	18
454	amazon	13	14
455	amazonie	13	14
456	amazonii	13	14
457	amazonky	13	14
458	amazonských	15	15
459	amazonu	13	14
460	ambasádě	18	18
461	ambasádu	16	16
462	ambasády	16	16
463	amber	13	14
464	ambice	13	14
465	ambicemi	13	14
466	ambici	13	14
467	ambicí	17	17
468	ambiciózní	19	19
469	ambivalentní	17	17
470	ambrož	15	15
471	ambulance	13	14
472	ambulanci	13	14
473	ambulantní	17	17
474	amd	13	14
475	amelia	13	14
476	amélie	18	18
477	amen	13	14
478	america	13	14
479	american	13	14
480	americe	13	14
481	americká	16	16
482	americké	18	18
483	amerického	18	18
484	americkém	18	18
485	americkému	18	18
486	americkou	13	14
487	americký	15	15
488	amerických	15	15
489	americkým	15	15
490	americkými	15	15
491	američan	17	17
492	američana	17	17
493	američané	18	18
494	američani	17	17
495	američanka	17	17
496	američanky	17	17
497	američanů	17	17
498	američanům	17	17
499	američany	17	17
500	američtí	17	17
501	amerika	13	14
502	amerikou	13	14
503	ameriku	13	14
504	ameriky	13	14
505	amfiteátru	16	16
506	aminokyselin	13	14
507	aminokyseliny	13	14
508	amnestie	13	14
509	amnestii	13	14
510	amoniaku	13	14
511	amplituda	13	14
512	amplitudy	13	14
513	amsterdam	13	14
514	amsterdamu	13	14
515	amsterodamu	13	14
516	amulet	13	14
517	amy	13	14
518	ana	12	12
519	anais	12	12
520	analogické	18	18
521	analogicky	13	12
522	analogie	12	12
523	analogii	12	12
524	analogií	17	17
525	analogové	18	18
526	analogového	18	18
527	analogový	15	15
528	analogových	15	15
529	analysis	13	12
530	analytici	13	12
531	analytická	16	16
532	analytické	18	18
533	analytického	18	18
534	analytický	15	15
535	analytických	15	15
536	analytik	13	12
537	analytika	13	12
538	analytiků	13	12
539	analýz	15	15
540	analýza	15	15
541	analýzách	16	16
542	analýze	15	15
543	analýzou	15	15
544	analyzoval	13	14
545	analyzovali	13	14
546	analyzovaných	15	15
547	analyzovat	13	14
548	analýzu	15	15
549	analyzuje	13	14
550	analyzují	17	17
551	analýzy	15	15
552	ananas	12	12
553	anarchie	12	12
554	anatomické	18	18
555	anatomicky	13	14
556	anatomických	15	15
557	anatomie	13	14
558	anatomii	13	14
559	anča	17	17
560	anděl	18	18
561	anděla	18	18
562	andělé	18	18
563	andělem	18	18
564	andělů	18	18
565	anděly	18	18
566	andílek	17	17
567	andreas	12	12
568	andrej	12	12
569	andreje	12	12
570	andrejem	13	14
571	android	12	12
572	aneb	12	12
573	anebo	12	12
574	anekdoty	13	12
575	anestezie	12	14
576	anestezii	12	14
577	aneta	12	12
578	anežka	15	15
579	anežky	15	15
580	angažmá	16	16
581	angažoval	15	15
582	angažování	17	17
583	angažovanost	15	15
584	angažovanosti	15	15
585	angažovat	15	15
586	angažuje	15	15
587	anglická	16	16
588	anglické	18	18
589	anglického	18	18
590	anglickém	18	18
591	anglickému	18	18
592	anglickou	12	12
593	anglicky	13	12
594	anglický	15	15
595	anglických	15	15
596	anglickým	15	15
597	anglickými	15	15
598	angličan	17	17
599	angličana	17	17
600	angličané	18	18
601	angličanka	17	17
602	angličanů	17	17
603	angličanům	17	17
604	angličany	17	17
605	angličtí	17	17
606	angličtina	17	17
607	angličtině	18	18
608	angličtinou	17	17
609	angličtinu	17	17
610	angličtiny	17	17
611	anglie	12	12
612	anglii	12	12
613	anglií	17	17
614	anička	17	17
615	aničku	17	17
616	aničky	17	17
617	animace	13	14
618	animované	18	18
619	animovaný	15	15
620	animovaných	15	15
621	anita	12	12
622	aniž	15	15
623	anketa	12	12
624	anketě	18	18
625	anketu	12	12
626	ankety	13	12
627	ano	12	12
628	anomálie	16	16
629	anonymitě	18	18
630	anonymity	13	14
631	anonymně	18	18
632	anonymní	17	17
633	anonymního	17	17
634	anonymních	17	17
635	anorexie	13	14
636	anorganických	15	15
637	antarktického	18	18
638	antarktidě	18	18
639	antarktidy	13	12
640	antén	18	18
641	anténa	18	18
642	anténu	18	18
643	antény	18	18
644	anthony	13	12
645	anthonyho	13	12
646	antibakteriální	17	17
647	antibiotik	12	12
648	antibiotika	12	12
649	antibiotiky	13	12
650	antické	18	18
651	antického	18	18
652	antických	15	15
653	antidepresiva	12	12
654	antikoncepce	12	12
655	antikoncepci	12	12
656	antikvariátu	16	16
657	antiky	13	12
658	antilopy	13	12
659	antimonopolní	17	17
660	antimonopolního	17	17
661	antioxidantů	13	14
662	antioxidanty	13	14
663	antisemitismu	13	14
664	antisemitismus	13	14
665	antoine	12	12
666	antologie	12	12
667	antologii	12	12
668	anton	12	12
669	antona	12	12
670	antonia	12	12
671	antonín	17	17
672	antonína	17	17
673	antonínem	17	17
674	antonínu	17	17
675	antonio	12	12
676	antropogenní	17	17
677	antropolog	12	12
678	antropologie	12	12
679	antropologii	12	12
680	apa	10	10
681	aparát	16	16
682	aparátem	16	16
683	aparátu	16	16
684	aparaturu	10	10
685	aparatury	13	10
686	apartmá	16	16
687	apartmánu	16	16
688	apartmány	16	16
689	apatie	10	10
690	apeluje	10	10
691	apetit	10	10
692	api	10	10
693	aplaus	10	10
694	aplikace	12	12
695	aplikacemi	13	14
696	aplikaci	12	12
697	aplikací	17	17
698	aplikacích	17	17
699	aplikacím	17	17
700	aplikační	17	17
701	aplikoval	11	11
702	aplikován	16	16
703	aplikována	16	16
704	aplikované	18	18
705	aplikovaného	18	18
706	aplikovaný	15	15
707	aplikovány	16	16
708	aplikovaných	15	15
709	aplikovat	11	11
710	aplikuje	10	10
711	aplikují	17	17
712	apod	10	10
713	apollo	10	10
714	apoštolů	18	18
715	app	10	10
716	apple	10	10
717	applu	10	10
718	apps	10	10
719	april	10	10
720	aps	10	10
721	aqua	10	10
722	aquaparku	10	10
723	arab	11	11
724	arábie	16	16
725	arábii	16	16
726	arabové	18	18
727	arabská	16	16
728	arabské	18	18
729	arabského	18	18
730	arabském	18	18
731	arabskou	11	11
732	arabsky	13	11
733	arabský	15	15
734	arabských	15	15
735	aranžmá	16	16
736	arbitráž	16	16
737	arcibiskup	12	12
738	arcibiskupa	12	12
739	arcibiskupem	13	14
740	arcibiskupství	17	17
741	areál	16	16
742	areálech	16	16
743	areálem	16	16
744	areálu	16	16
745	areálů	16	16
746	areály	16	16
747	arena	12	12
748	aréna	18	18
749	areně	18	18
750	aréně	18	18
751	areny	13	12
752	arény	18	18
753	argentina	12	12
754	argentině	18	18
755	argentinské	18	18
756	argentiny	13	12
757	argument	13	14
758	argumentace	13	14
759	argumentaci	13	14
760	argumentem	13	14
761	argumentoval	13	14
762	argumentovat	13	14
763	argumentu	13	14
764	argumentů	13	14
765	argumentuje	13	14
766	argumentují	17	17
767	argumenty	13	14
768	archaické	18	18
769	archeolog	12	12
770	archeologické	18	18
771	archeologického	18	18
772	archeologický	15	15
773	archeologických	15	15
774	archeologie	12	12
775	archeologii	12	12
776	archeologové	18	18
777	archeologů	12	12
778	archetypální	17	17
779	architekta	12	12
780	architektem	13	14
781	architekti	12	12
782	architektka	12	12
783	architektonická	16	16
784	architektonické	18	18
785	architektonického	18	18
786	architektonickou	12	12
787	architektonicky	13	12
788	architektonický	15	15
789	architektonických	15	15
790	architektonickým	15	15
791	architektů	12	12
792	architektům	13	14
793	architektura	12	12
794	architekturou	12	12
795	architekturu	12	12
796	architektury	13	12
797	architektuře	16	16
798	architekty	13	12
799	archiv	12	12
800	archivace	12	12
801	archivaci	12	12
802	archivech	12	12
803	archivní	17	17
804	archivních	17	17
805	archivu	12	12
806	archivů	12	12
807	archivy	13	12
808	archy	13	12
809	ariel	8	8
810	aristokracie	12	12
811	aristotela	9	9
812	aristoteles	9	9
813	aristotelés	18	18
814	arizoně	18	18
815	arktidě	18	18
816	armád	16	16
817	armáda	16	16
818	armádě	18	18
819	armádní	17	17
820	armádního	17	17
821	armádních	17	17
822	armádou	16	16
823	armádu	16	16
824	armády	16	16
825	armand	13	14
826	armanda	13	14
827	armani	13	14
828	arménie	18	18
829	armstrong	13	14
830	armstronga	13	14
831	arne	12	12
832	arneho	12	12
833	arnold	12	12
834	arnolda	12	12
835	arnošt	18	18
836	arnošta	18	18
837	aro	9	9
838	arogance	12	12
839	aroganci	12	12
840	arogantně	18	18
841	arogantní	17	17
842	aroma	13	14
843	aromatické	18	18
844	aromatických	15	15
845	arpád	16	16
846	arriane	12	12
847	arsenal	12	12
848	arsenalu	12	12
849	art	7	7
850	artefakt	8	8
851	artefaktů	8	8
852	artefakty	13	8
853	arthur	7	7
854	arthura	7	7
855	artiklem	13	14
856	artis	8	8
857	artritida	8	8
858	artritidou	9	9
859	artritidy	13	8
860	artur	7	7
861	artura	7	7
862	arzenál	16	16
863	arzenálu	16	16
864	asfalt	5	6
865	asfaltové	18	18
866	asfaltový	15	15
867	asfaltových	15	15
868	ashley	13	8
869	asi	8	8
870	asie	8	8
871	asijské	18	18
872	asijského	18	18
873	asijský	15	15
874	asijských	15	15
875	asistence	12	12
876	asistenci	12	12
877	asistencí	17	17
878	asistenční	17	17
879	asistent	12	12
880	asistenta	12	12
881	asistentek	12	12
882	asistentem	13	14
883	asistenti	12	12
884	asistentka	12	12
885	asistentku	12	12
886	asistentky	13	12
887	asistentů	12	12
888	asistenty	13	12
889	asistoval	11	11
890	asistované	18	18
891	asociace	12	12
892	asociaci	12	12
893	asociací	17	17
894	aspekt	10	10
895	aspektech	12	12
896	aspektem	13	14
897	aspektu	10	10
898	aspektů	10	10
899	aspektům	13	14
900	aspekty	13	10
901	aspire	10	10
902	aspirin	12	12
903	aspoň	19	19
904	astma	13	14
905	astmatu	13	14
906	aston	12	12
907	astra	7	7
908	astrid	8	8
909	astronaut	12	12
910	astronom	13	14
911	astronomické	18	18
912	astronomických	15	15
913	astronomie	13	14
914	astronomii	13	14
915	astronomové	18	18
916	astronomů	13	14
917	asus	7	7
918	asymetrické	18	18
919	asymetrie	13	14
920	ať	19	19
921	atari	8	8
922	atd	5	6
923	ateliér	18	18
924	ateliéru	18	18
925	ateliérů	18	18
926	ateliéry	18	18
927	atentát	16	16
928	atentátu	16	16
929	atény	18	18
930	athénách	18	18
931	atlantě	18	18
932	atlantic	12	12
933	atlantik	12	12
934	atlantiku	12	12
935	atlas	5	6
936	atlasu	7	7
937	atlet	8	8
938	atleti	8	8
939	atletice	12	12
940	atletické	18	18
941	atletického	18	18
942	atletický	15	15
943	atletických	15	15
944	atletika	8	8
945	atletiky	13	8
946	atletů	8	8
947	atmosféra	18	18
948	atmosférické	18	18
949	atmosférou	18	18
950	atmosféru	18	18
951	atmosféry	18	18
952	atmosféře	18	18
953	atom	13	14
954	atomová	16	16
955	atomové	18	18
956	atomového	18	18
957	atomovou	13	14
958	atomových	15	15
959	atomu	13	14
960	atomů	13	14
961	atomy	13	14
962	atrakce	12	12
963	atrakci	12	12
964	atrakcí	17	17
965	atraktivitu	11	11
966	atraktivity	13	11
967	atraktivnější	18	18
968	atraktivní	17	17
969	atraktivního	17	17
970	atraktivních	17	17
971	atraktivním	17	17
972	atraktoru	9	9
973	atribut	11	11
974	atributů	11	11
975	atributy	13	11
976	atypické	18	18
977	au	7	7
978	audi	8	8
979	audie	8	8
980	audieho	9	9
981	audience	12	12
982	audienci	12	12
983	audio	9	9
984	audiovizuální	17	17
985	audit	8	8
986	auditor	9	9
987	auditorů	9	9
988	auditoři	16	16
989	august	7	7
990	augusta	7	7
991	augustin	12	12
992	augustina	12	12
993	aukce	12	12
994	aukci	12	12
995	aukcí	17	17
996	aukcích	17	17
997	aukční	17	17
998	aura	7	7
999	australan	12	12
1000	austrálie	16	16
1001	austrálii	16	16
1002	australská	16	16
1003	australské	18	18
1004	australského	18	18
1005	australský	15	15
1006	australských	15	15
1007	aut	7	7
1008	auta	7	7
1009	autě	18	18
1010	autech	12	12
1011	autem	13	14
1012	autentická	16	16
1013	autentické	18	18
1014	autentickou	12	12
1015	autenticky	13	12
1016	autentický	15	15
1017	autentických	15	15
1018	autíčka	17	17
1019	autíčko	17	17
1020	autismem	13	14
1021	auto	9	9
1022	autobiografie	11	11
1023	autobiografii	11	11
1024	autobus	11	11
1025	autobuse	11	11
1026	autobusech	12	12
1027	autobusem	13	14
1028	autobusová	16	16
1029	autobusové	18	18
1030	autobusového	18	18
1031	autobusovou	11	11
1032	autobusových	15	15
1033	autobusu	11	11
1034	autobusů	11	11
1035	autobusy	13	11
1036	automat	13	14
1037	automatech	13	14
1038	automatem	13	14
1039	automatická	16	16
1040	automatické	18	18
1041	automatického	18	18
1042	automatickou	13	14
1043	automaticky	13	14
1044	automatický	15	15
1045	automatických	15	15
1046	automatickým	15	15
1047	automatizace	13	14
1048	automatizaci	13	14
1049	automatizované	18	18
1050	automatu	13	14
1051	automatů	13	14
1052	automaty	13	14
1053	automobil	13	14
1054	automobilce	13	14
1055	automobilech	13	14
1056	automobilek	13	14
1057	automobilem	13	14
1058	automobilka	13	14
1059	automobilku	13	14
1060	automobilky	13	14
1061	automobilová	16	16
1062	automobilové	18	18
1063	automobilového	18	18
1064	automobilovém	18	18
1065	automobilovou	13	14
1066	automobilový	15	15
1067	automobilových	15	15
1068	automobilu	13	14
1069	automobilů	13	14
1070	automobily	13	14
1071	automotive	13	14
1072	autonehodě	18	18
1073	autonomie	13	14
1074	autonomii	13	14
1075	autonomní	17	17
1076	autoportrét	18	18
1077	autor	9	9
1078	autora	9	9
1079	autorem	13	14
1080	autorit	9	9
1081	autorita	9	9
1082	autoritám	16	16
1083	autoritativně	18	18
1084	autoritativní	17	17
1085	autoritě	18	18
1086	autoritou	9	9
1087	autoritu	9	9
1088	autority	13	9
1089	autorizovaných	15	15
1090	autorka	9	9
1091	autorkou	9	9
1092	autorky	13	9
1093	autorova	11	11
1094	autorově	18	18
1095	autorovi	11	11
1096	autorovy	13	11
1097	autorská	16	16
1098	autorské	18	18
1099	autorského	18	18
1100	autorskou	9	9
1101	autorský	15	15
1102	autorských	15	15
1103	autorským	15	15
1104	autorství	17	17
1105	autorů	9	9
1106	autorům	13	14
1107	autorův	11	11
1108	autory	13	9
1109	autoři	16	16
1110	autosalonu	12	12
1111	autu	7	7
1112	autům	13	14
1113	auty	13	7
1114	avantgardní	17	17
1115	avatar	11	11
1116	avizoval	11	14
1117	avšak	18	18
1118	axel	13	14
1119	azyl	13	14
1120	až	15	15
1121	bába	16	16
1122	babek	11	11
1123	babi	11	11
1124	babičce	17	17
1125	babiččině	18	18
1126	babiččiny	17	17
1127	babiček	17	17
1128	babička	17	17
1129	babičko	17	17
1130	babičkou	17	17
1131	babičku	17	17
1132	babičky	17	17
1133	babinského	18	18
1134	babinský	15	15
1135	babiš	18	18
1136	babiše	18	18
1137	babišem	18	18
1138	babišovi	18	18
1139	babka	11	11
1140	babku	11	11
1141	babky	13	11
1142	bábovku	16	16
1143	bábu	16	16
1144	baby	13	11
1145	báby	16	16
1146	babylon	13	12
1147	back	12	12
1148	bačkory	17	17
1149	bad	11	11
1150	bádání	17	17
1151	badatel	11	11
1152	badatele	11	11
1153	badatelé	18	18
1154	badatelské	18	18
1155	badatelů	11	11
1156	badminton	13	14
1157	bagdádu	16	16
1158	bagetu	11	11
1159	bagety	13	11
1160	bagr	11	11
1161	bahenní	17	17
1162	bahna	12	12
1163	bahně	18	18
1164	bahnem	13	14
1165	bahno	12	12
1166	bacha	12	12
1167	báječná	17	17
1168	báječné	18	18
1169	báječně	18	18
1170	báječnou	17	17
1171	báječný	17	17
1172	bájné	18	18
1173	bakalářské	18	18
1174	bakteriální	17	17
1175	bakteriálních	17	17
1176	bakterie	11	11
1177	bakteriemi	13	14
1178	bakterii	11	11
1179	bakterií	17	17
1180	bakteriím	17	17
1181	bál	16	16
1182	bála	16	16
1183	balady	13	11
1184	balancuje	12	12
1185	balené	18	18
1186	balení	17	17
1187	balet	11	11
1188	baletní	17	17
1189	baletu	11	11
1190	bali	11	11
1191	balí	17	17
1192	balíček	17	17
1193	balíčkem	17	17
1194	balíčku	17	17
1195	balíčků	17	17
1196	balíčky	17	17
1197	balík	17	17
1198	balíkem	17	17
1199	balíku	17	17
1200	balíků	17	17
1201	balíky	17	17
1202	balil	11	11
1203	balit	11	11
1204	balkán	16	16
1205	balkáně	18	18
1206	balkánu	16	16
1207	balkon	12	12
1208	balkón	19	19
1209	balkoně	18	18
1210	balkonu	12	12
1211	balkonů	12	12
1212	balkony	13	12
1213	balon	12	12
1214	balón	19	19
1215	balonek	12	12
1216	balonem	13	14
1217	balonky	13	12
1218	balonu	12	12
1219	balony	13	12
1220	baltského	18	18
1221	baltu	11	11
1222	balvan	12	12
1223	balvanu	12	12
1224	balvanů	12	12
1225	balvany	13	12
1226	balzám	16	16
1227	bambus	13	14
1228	bambusové	18	18
1229	bambusu	13	14
1230	banality	13	12
1231	banálně	18	18
1232	banální	17	17
1233	banán	16	16
1234	banánů	16	16
1235	banány	16	16
1236	banda	12	12
1237	banka	12	12
1238	bankách	16	16
1239	bankám	16	16
1240	bankami	13	14
1241	bankéř	18	18
1242	bankéře	18	18
1243	bankéři	18	18
1244	bankéřů	18	18
1245	bankomatu	13	14
1246	bankomatů	13	14
1247	bankou	12	12
1248	bankovek	12	12
1249	bankovku	12	12
1250	bankovky	13	12
1251	bankovní	17	17
1252	bankovnictví	17	17
1253	bankovního	17	17
1254	bankovních	17	17
1255	bankovním	17	17
1256	bankrot	12	12
1257	bankrotu	12	12
1258	banku	12	12
1259	banky	13	12
1260	baňky	19	19
1261	banské	18	18
1262	bar	11	11
1263	bára	16	16
1264	barák	16	16
1265	barákem	16	16
1266	baráku	16	16
1267	baráků	16	16
1268	baráky	16	16
1269	barbara	11	11
1270	barbarské	18	18
1271	barbarství	17	17
1272	barbarů	11	11
1273	barbary	13	11
1274	barbaře	16	16
1275	barbaři	16	16
1276	barbie	11	11
1277	barbora	11	11
1278	barbory	13	11
1279	barcelona	12	12
1280	barceloně	18	18
1281	barcelonou	12	12
1282	barcelonu	12	12
1283	barcelony	13	12
1284	barech	12	12
1285	barel	11	11
1286	barelů	11	11
1287	barem	13	14
1288	bareš	18	18
1289	baret	11	11
1290	barev	11	11
1291	barevná	16	16
1292	barevné	18	18
1293	barevně	18	18
1294	barevného	18	18
1295	barevném	18	18
1296	barevnost	12	12
1297	barevnosti	12	12
1298	barevností	17	17
1299	barevnou	12	12
1300	barevný	15	15
1301	barevných	15	15
1302	barevným	15	15
1303	barevnými	15	15
1304	bariér	18	18
1305	bariéra	18	18
1306	bariérou	18	18
1307	bariéru	18	18
1308	bariéry	18	18
1309	barikády	16	16
1310	barman	13	14
1311	barmana	13	14
1312	barmě	18	18
1313	barnabáš	18	18
1314	baroka	11	11
1315	barokní	17	17
1316	barokního	17	17
1317	barokních	17	17
1318	barokním	17	17
1319	baroko	11	11
1320	baroku	11	11
1321	baron	12	12
1322	barona	12	12
1323	baronka	12	12
1324	bárou	16	16
1325	barrandov	12	12
1326	bárta	16	16
1327	barták	16	16
1328	bartáka	16	16
1329	bartoloměje	18	18
1330	bartoš	18	18
1331	bartoška	18	18
1332	bartošová	18	18
1333	bartošové	18	18
1334	bártou	16	16
1335	bártu	16	16
1336	bárty	16	16
1337	baru	11	11
1338	barů	11	11
1339	báru	16	16
1340	barva	11	11
1341	barvách	16	16
1342	barvám	16	16
1343	barvami	13	14
1344	barvě	18	18
1345	barvení	17	17
1346	barví	17	17
1347	barvitě	18	18
1348	barviv	11	11
1349	barviva	11	11
1350	barvivo	11	11
1351	barvou	11	11
1352	barvu	11	11
1353	barvy	13	11
1354	bary	13	11
1355	báře	16	16
1356	baseball	11	11
1357	baseballové	18	18
1358	baseballovou	11	11
1359	baseballu	11	11
1360	báseň	19	19
1361	basket	11	11
1362	basketbal	11	11
1363	basketbalista	11	11
1364	basketbalisté	18	18
1365	basketbalové	18	18
1366	basketbalu	11	11
1367	básně	18	18
1368	básni	16	16
1369	básní	17	17
1370	básníci	17	17
1371	básnické	18	18
1372	básnickou	16	16
1373	básnických	16	16
1374	básnictví	17	17
1375	básničku	17	17
1376	básničky	17	17
1377	básních	17	17
1378	básník	17	17
1379	básníka	17	17
1380	básníkem	17	17
1381	básníkovi	17	17
1382	básníků	17	17
1383	básníky	17	17
1384	básnířka	17	17
1385	basta	11	11
1386	basu	11	11
1387	basů	11	11
1388	basy	13	11
1389	bašty	18	18
1390	bát	16	16
1391	baťa	19	19
1392	baterie	11	11
1393	bateriemi	13	14
1394	baterii	11	11
1395	baterií	17	17
1396	baterkou	11	11
1397	baterku	11	11
1398	baterky	13	11
1399	batman	13	14
1400	batoh	11	11
1401	batohem	13	14
1402	batohu	11	11
1403	batohy	13	11
1404	batolata	11	11
1405	batole	11	11
1406	baví	17	17
1407	bavič	17	17
1408	bavil	11	11
1409	bavila	11	11
1410	bavili	11	11
1411	bavilo	11	11
1412	bavily	13	11
1413	bavím	17	17
1414	bavíme	17	17
1415	bavit	11	11
1416	bavlna	12	12
1417	bavlněné	18	18
1418	bavlny	13	12
1419	bavorska	11	11
1420	bavorské	18	18
1421	bavorského	18	18
1422	bavorsko	11	11
1423	bavorsku	11	11
1424	bayern	13	12
1425	bayernu	13	12
1426	bazalech	12	14
1427	bazalky	13	14
1428	bazální	17	17
1429	bazálních	17	17
1430	bazar	11	14
1431	bazaru	11	14
1432	báze	16	16
1433	bazén	18	18
1434	bázeň	19	19
1435	bazénem	18	18
1436	bazénové	18	18
1437	bazénu	18	18
1438	bazénů	18	18
1439	bazény	18	18
1440	bázi	16	16
1441	bází	17	17
1442	bazických	15	15
1443	bazilice	12	14
1444	baziliky	13	14
1445	bázlivě	18	18
1446	bažin	15	15
1447	bažiny	15	15
1448	bbc	12	12
1449	bdělosti	18	18
1450	bdění	18	18
1451	bdí	17	17
1452	bea	11	11
1453	beach	12	12
1454	bean	12	12
1455	beatles	11	11
1456	beatrice	12	12
1457	beaufort	11	11
1458	beck	12	12
1459	beckham	13	14
1460	becky	13	12
1461	béčka	18	18
1462	béčko	18	18
1463	bečva	17	17
1464	bečvou	17	17
1465	bečvy	17	17
1466	běda	18	18
1467	beden	12	12
1468	bederní	17	17
1469	bedlivě	18	18
1470	bedna	12	12
1471	bednář	16	16
1472	bedně	18	18
1473	bednění	18	18
1474	bednu	12	12
1475	bedny	13	12
1476	bedra	11	11
1477	bedrech	12	12
1478	bedřich	16	16
1479	bedřicha	16	16
1480	bedýnky	15	15
1481	bees	11	11
1482	beethoven	12	12
1483	begley	13	11
1484	běh	18	18
1485	běhá	18	18
1486	běhají	18	18
1487	běhal	18	18
1488	běhala	18	18
1489	běhali	18	18
1490	běhaly	18	18
1491	běhání	18	18
1492	běhat	18	18
1493	behaviorální	17	17
1494	během	18	18
1495	běhounek	18	18
1496	běhu	18	18
1497	běhy	18	18
1498	bechyně	18	18
1499	bejrútu	11	11
1500	bejt	11	11
1501	bek	11	11
1502	beka	11	11
1503	běla	18	18
1504	bělavé	18	18
1505	bělé	18	18
1506	bělehrad	18	18
1507	bělehradě	18	18
1508	bělehradu	18	18
1509	beletrie	11	11
1510	beletrii	11	11
1511	belgické	18	18
1512	belgického	18	18
1513	belgický	15	15
1514	belgie	11	11
1515	belgii	11	11
1516	bell	11	11
1517	bella	11	11
1518	belle	11	11
1519	bello	11	11
1520	bělobrádek	18	18
1521	běloch	18	18
1522	bělochů	18	18
1523	bělochy	18	18
1524	běloruska	18	18
1525	bělorusko	18	18
1526	bělorusku	18	18
1527	bělostné	18	18
1528	běloši	18	18
1529	bém	18	18
1530	béma	18	18
1531	ben	12	12
1532	bena	12	12
1533	benátek	16	16
1534	benátkách	16	16
1535	benátky	16	16
1536	benda	12	12
1537	bendl	12	12
1538	bene	12	12
1539	benedikt	12	12
1540	benedikta	12	12
1541	benefiční	17	17
1542	benefitů	12	12
1543	benefity	13	12
1544	benem	13	14
1545	beneš	18	18
1546	beneše	18	18
1547	benešov	18	18
1548	benešova	18	18
1549	benešová	18	18
1550	benešové	18	18
1551	benešově	18	18
1552	benešovi	18	18
1553	benghází	17	17
1554	benjamin	13	14
1555	benjamina	13	14
1556	benno	12	12
1557	benny	13	12
1558	benovi	12	12
1559	bentley	13	12
1560	benzin	12	14
1561	benzín	17	17
1562	benzinem	13	14
1563	benzinové	18	18
1564	benzínové	18	18
1565	benzinový	15	15
1566	benzinových	15	15
1567	benzinu	12	14
1568	benzínu	17	17
1569	ber	11	11
1570	beran	12	12
1571	berana	12	12
1572	beránek	16	16
1573	beránka	16	16
1574	berdych	13	12
1575	berdycha	13	12
1576	bere	11	11
1577	bereme	13	14
1578	berenika	12	12
1579	bereníké	18	18
1580	bereš	18	18
1581	berete	11	11
1582	berg	11	11
1583	berger	11	11
1584	bergmann	13	14
1585	berija	11	11
1586	berit	11	11
1587	berka	11	11
1588	berkeley	13	11
1589	berko	11	11
1590	berle	11	11
1591	berlích	17	17
1592	berlín	17	17
1593	berlína	17	17
1594	berlíně	18	18
1595	berlínem	17	17
1596	berlínské	18	18
1597	berlínského	18	18
1598	berlínském	18	18
1599	berlusconi	12	12
1600	berlusconiho	12	12
1601	bernard	12	12
1602	bernarda	12	12
1603	bernie	12	12
1604	bernu	12	12
1605	berou	11	11
1606	beroun	12	12
1607	berouna	12	12
1608	berouně	18	18
1609	berounky	13	12
1610	bert	11	11
1611	berta	11	11
1612	berte	11	11
1613	beru	11	11
1614	beryl	13	11
1615	beseda	11	11
1616	besedě	18	18
1617	besedu	11	11
1618	besedy	13	11
1619	besian	12	12
1620	beskyd	13	11
1621	beskydech	13	12
1622	beskydy	13	11
1623	běsnění	18	18
1624	bess	11	11
1625	besser	11	11
1626	best	11	11
1627	bestie	11	11
1628	bestii	11	11
1629	bestseller	11	11
1630	bestsellerem	13	14
1631	bestselleru	11	11
1632	beta	11	11
1633	beth	11	11
1634	betlém	18	18
1635	betléma	18	18
1636	betlémské	18	18
1637	betlémů	18	18
1638	betlémy	18	18
1639	beton	12	12
1640	betonem	13	14
1641	betonová	16	16
1642	betonové	18	18
1643	betonového	18	18
1644	betonovém	18	18
1645	betonovou	12	12
1646	betonový	15	15
1647	betonových	15	15
1648	betonovými	15	15
1649	betonu	12	12
1650	betonů	12	12
1651	betony	13	12
1652	betsy	13	11
1653	betty	13	11
1654	beverly	13	11
1655	bez	11	14
1656	bezbariérové	18	18
1657	bezbariérový	18	18
1658	bezbarvé	18	18
1659	bezbarvý	15	15
1660	bezbarvým	15	15
1661	bezbolestně	18	18
1662	bezbranná	16	16
1663	bezbranné	18	18
1664	bezcenné	18	18
1665	bezcílně	18	18
1666	bezčasí	17	17
1667	bezděčně	18	18
1668	bezděky	18	18
1669	bezdomovce	13	14
1670	bezdomovci	13	14
1671	bezdomovců	13	14
1672	bezdomovec	13	14
1673	bezdrátová	16	16
1674	bezdrátové	18	18
1675	bezdrátově	18	18
1676	bezdrátového	18	18
1677	bezdrátovou	16	16
1678	bezdrátový	16	16
1679	bezdrátových	16	16
1680	bezdůvodně	18	18
1681	beze	11	14
1682	bezedné	18	18
1683	bezelstně	18	18
1684	bezesporu	11	14
1685	bezezbytku	13	14
1686	bezhlavě	18	18
1687	bezchybné	18	18
1688	bezchybně	18	18
1689	bezchybný	15	15
1690	bezkonkurenčně	18	18
1691	bezkontaktní	17	17
1692	bezmála	16	16
1693	bezmezně	18	18
1694	bezmoc	13	14
1695	bezmoci	13	14
1696	bezmocná	16	16
1697	bezmocné	18	18
1698	bezmocně	18	18
1699	bezmocní	17	17
1700	bezmocný	15	15
1701	bezmyšlenkovitě	18	18
1702	beznaděj	18	18
1703	beznaděje	18	18
1704	beznaději	18	18
1705	beznadějné	18	18
1706	beznadějně	18	18
1707	beznadějný	18	18
1708	bezobratlých	15	15
1709	bezodkladně	18	18
1710	bezohledně	18	18
1711	bezohlednost	12	14
1712	bezohledný	15	15
1713	bezostyšně	18	18
1714	bezpečí	17	17
1715	bezpečná	17	17
1716	bezpečné	18	18
1717	bezpečně	18	18
1718	bezpečného	18	18
1719	bezpečněji	18	18
1720	bezpečnější	18	18
1721	bezpečném	18	18
1722	bezpečnost	17	17
1723	bezpečnosti	17	17
1724	bezpečností	17	17
1725	bezpečnostní	17	17
1726	bezpečnostního	17	17
1727	bezpečnostních	17	17
1728	bezpečnostním	17	17
1729	bezpečnostními	17	17
1730	bezpečnou	17	17
1731	bezpečný	17	17
1732	bezpečných	17	17
1733	bezpilotní	17	17
1734	bezplatná	16	16
1735	bezplatné	18	18
1736	bezplatně	18	18
1737	bezplatnou	12	14
1738	bezplatný	15	15
1739	bezpočet	17	17
1740	bezpočtu	17	17
1741	bezpodmínečně	18	18
1742	bezpochyby	13	14
1743	bezpráví	17	17
1744	bezprecedentní	17	17
1745	bezproblémová	18	18
1746	bezproblémové	18	18
1747	bezproblémový	18	18
1748	bezprostředně	18	18
1749	bezprostřední	17	17
1750	bezprostředního	17	17
1751	bezprostředním	17	17
1752	bezradně	18	18
1753	bezradnost	12	14
1754	bezradný	15	15
1755	bezstarostné	18	18
1756	bezstarostně	18	18
1757	bezstarostný	15	15
1758	beztak	11	14
1759	beztíže	17	17
1760	beztrestně	18	18
1761	bezu	11	14
1762	bezva	11	14
1763	bezvadné	18	18
1764	bezvadně	18	18
1765	bezvadný	15	15
1766	bezvědomí	18	18
1767	bezvětří	18	18
1768	bezvládné	18	18
1769	bezvládně	18	18
1770	bezvýhradně	18	18
1771	bezvýsledně	18	18
1772	bezvýznamná	16	16
1773	bezvýznamné	18	18
1774	bezvýznamný	15	15
1775	bezvýznamných	15	15
1776	bezzubka	11	14
1777	běž	18	18
1778	běžce	18	18
1779	běžci	18	18
1780	běžců	18	18
1781	běžec	18	18
1782	běžecké	18	18
1783	běžeckého	18	18
1784	běžeckých	18	18
1785	běžel	18	18
1786	běžela	18	18
1787	běželi	18	18
1788	běželo	18	18
1789	běžely	18	18
1790	běžet	18	18
1791	běží	18	18
1792	běžící	18	18
1793	běžícím	18	18
1794	běžím	18	18
1795	běžkách	18	18
1796	běžky	18	18
1797	běžná	18	18
1798	běžné	18	18
1799	běžně	18	18
1800	běžného	18	18
1801	běžnější	18	18
1802	běžném	18	18
1803	běžnému	18	18
1804	běžní	18	18
1805	běžnou	18	18
1806	běžný	18	18
1807	běžných	18	18
1808	běžným	18	18
1809	běžnými	18	18
1810	béžové	18	18
1811	běžte	18	18
1812	bianca	12	12
1813	bibi	11	11
1814	bible	11	11
1815	bibli	11	11
1816	biblické	18	18
1817	biblického	18	18
1818	biblický	15	15
1819	biblických	15	15
1820	bicí	17	17
1821	bicích	17	17
1822	bicykl	13	12
1823	bicyklu	13	12
1824	bicykly	13	12
1825	bič	17	17
1826	biče	17	17
1827	bičem	17	17
1828	bída	17	17
1829	bídě	18	18
1830	bídně	18	18
1831	bídou	17	17
1832	bídu	17	17
1833	bídy	17	17
1834	bienále	16	16
1835	big	11	11
1836	biggar	11	11
1837	bije	11	11
1838	bijí	17	17
1839	bil	11	11
1840	bila	11	11
1841	bílá	17	17
1842	bilance	12	12
1843	bilanci	12	12
1844	bilancí	17	17
1845	bild	11	11
1846	bíle	17	17
1847	bílé	18	18
1848	bílého	18	18
1849	bílek	17	17
1850	bílém	18	18
1851	bílému	18	18
1852	bili	11	11
1853	bílí	17	17
1854	bílina	17	17
1855	bíliny	17	17
1856	bilion	12	12
1857	bilionu	12	12
1858	bilionů	12	12
1859	bílka	17	17
1860	bílkovin	17	17
1861	bílkoviny	17	17
1862	bílků	17	17
1863	bílky	17	17
1864	bill	11	11
1865	billa	11	11
1866	billboard	11	11
1867	billboardů	11	11
1868	billboardy	13	11
1869	bille	11	11
1870	billem	13	14
1871	billovi	11	11
1872	billy	13	11
1873	billyho	13	11
1874	bilo	11	11
1875	bílo	17	17
1876	bílou	17	17
1877	bílý	17	17
1878	bílých	17	17
1879	bílým	17	17
1880	bílými	17	17
1881	bim	13	14
1882	bin	12	12
1883	bina	12	12
1884	binární	17	17
1885	binec	12	12
1886	bio	11	11
1887	biodiverzity	13	14
1888	biografie	11	11
1889	biografii	11	11
1890	biografu	11	11
1891	biochemické	18	18
1892	biochemie	13	14
1893	biolog	11	11
1894	biologická	16	16
1895	biologické	18	18
1896	biologického	18	18
1897	biologickou	12	12
1898	biologicky	13	12
1899	biologický	15	15
1900	biologických	15	15
1901	biologickým	15	15
1902	biologickými	15	15
1903	biologie	11	11
1904	biologii	11	11
1905	biologové	18	18
1906	biologů	11	11
1907	biomasa	13	14
1908	biomasu	13	14
1909	biomasy	13	14
1910	bionomie	13	14
1911	biopaliv	11	11
1912	biopaliva	11	11
1913	bioplynové	18	18
1914	bioplynových	15	15
1915	bioplynu	13	12
1916	biopotravin	12	12
1917	biopotraviny	13	12
1918	biotop	11	11
1919	biotopu	11	11
1920	biotopů	11	11
1921	biotopy	13	11
1922	bipolární	17	17
1923	bipolárního	17	17
1924	birminghamu	13	14
1925	bis	11	11
1926	biskup	11	11
1927	biskupa	11	11
1928	biskupem	13	14
1929	biskupové	18	18
1930	biskupovi	11	11
1931	biskupské	18	18
1932	biskupství	17	17
1933	biskupů	11	11
1934	biskupy	13	11
1935	bistra	11	11
1936	bistro	11	11
1937	bistru	11	11
1938	bit	11	11
1939	bít	17	17
1940	bitev	11	11
1941	bitevní	17	17
1942	bitevním	17	17
1943	bití	17	17
1944	bitky	13	11
1945	bitů	11	11
1946	bitva	11	11
1947	bitvách	16	16
1948	bitvě	18	18
1949	bitvou	11	11
1950	bitvu	11	11
1951	bitvy	13	11
1952	bivoltinní	17	17
1953	bizarní	17	17
1954	bizarních	17	17
1955	bižuterie	15	15
1956	bjørne	-1	-1
1957	bla	11	11
1958	black	12	12
1959	blackberry	13	12
1960	blackhorse	12	12
1961	blaha	11	11
1962	bláha	16	16
1963	blahem	13	14
1964	blaho	11	11
1965	blahobyt	13	11
1966	blahobytu	13	11
1967	blahodárné	18	18
1968	blahodárně	18	18
1969	blahodárný	16	16
1970	blahopřání	17	17
1971	blahosklonně	18	18
1972	bláhové	18	18
1973	blair	11	11
1974	blaira	11	11
1975	blake	11	11
1976	blanc	12	12
1977	blanensku	12	12
1978	blanche	12	12
1979	blaník	17	17
1980	blank	12	12
1981	blanka	12	12
1982	blanky	13	12
1983	blankytně	18	18
1984	blanska	12	12
1985	blansko	12	12
1986	blansku	12	12
1987	bláta	16	16
1988	blátě	18	18
1989	blátem	16	16
1990	blatné	18	18
1991	bláto	16	16
1992	blaze	11	14
1993	blázen	16	16
1994	blázince	16	16
1995	blázinci	16	16
1996	blázinec	16	16
1997	blázna	16	16
1998	blázni	16	16
1999	bláznivá	16	16
2000	bláznivé	18	18
2001	bláznivě	18	18
2002	bláznivého	18	18
2003	bláznivou	16	16
2004	bláznivý	16	16
2005	bláznivých	16	16
2006	bláznovství	17	17
2007	bláznů	16	16
2008	blázny	16	16
2009	blažek	15	15
2010	blažené	18	18
2011	blaženě	18	18
2012	blaženosti	15	15
2013	blažený	15	15
2014	blažka	15	15
2015	blbá	16	16
2016	blbce	12	12
2017	blbci	12	12
2018	blbé	18	18
2019	blbě	18	18
2020	blbec	12	12
2021	blbej	11	11
2022	blbost	11	11
2023	blbosti	11	11
2024	blbou	11	11
2025	blbý	15	15
2026	bledá	16	16
2027	bledé	18	18
2028	bledě	18	18
2029	bledého	18	18
2030	bledne	12	12
2031	bledou	11	11
2032	bledší	18	18
2033	bledý	15	15
2034	bledých	15	15
2035	blechy	13	12
2036	blesk	11	11
2037	bleskem	13	14
2038	blesklo	11	11
2039	bleskové	18	18
2040	bleskově	18	18
2041	blesku	11	11
2042	blesků	11	11
2043	bleskurychle	13	12
2044	blesky	13	11
2045	bliká	16	16
2046	blikající	17	17
2047	blízce	17	17
2048	blízcí	17	17
2049	blízká	17	17
2050	blízké	18	18
2051	blízkého	18	18
2052	blízkém	18	18
2053	blízko	17	17
2054	blízkost	17	17
2055	blízkosti	17	17
2056	blízkou	17	17
2057	blízký	17	17
2058	blízkých	17	17
2059	blízkým	17	17
2060	blízkými	17	17
2061	blíž	17	17
2062	blíže	17	17
2063	blíží	17	17
2064	blížící	17	17
2065	blížícího	17	17
2066	blížících	17	17
2067	blížícím	17	17
2068	blížil	17	17
2069	blížila	17	17
2070	blížili	17	17
2071	blížilo	17	17
2072	blížily	17	17
2073	blížíme	17	17
2074	blížit	17	17
2075	bližní	17	17
2076	bližního	17	17
2077	bližních	17	17
2078	bližním	17	17
2079	bližnímu	17	17
2080	bližší	18	18
2081	bližšího	18	18
2082	bližším	18	18
2083	blocích	17	17
2084	blog	11	11
2085	blogu	11	11
2086	blogů	11	11
2087	blogy	13	11
2088	blok	11	11
2089	blokády	16	16
2090	blokem	13	14
2091	blokování	17	17
2092	blokovat	11	11
2093	bloku	11	11
2094	bloků	11	11
2095	blokuje	11	11
2096	blokují	17	17
2097	bloky	13	11
2098	blond	12	12
2099	blonďák	19	19
2100	blonďatá	19	19
2101	blonďaté	19	19
2102	blondýna	15	15
2103	blondýnka	15	15
2104	blondýny	15	15
2105	bloomberg	13	14
2106	bloudění	18	18
2107	bloudí	17	17
2108	bloudil	11	11
2109	bloudit	11	11
2110	bludiště	18	18
2111	bludišti	18	18
2112	bludy	13	11
2113	blue	11	11
2114	blues	11	11
2115	bluetooth	11	11
2116	blůzu	11	14
2117	blůzy	13	14
2118	blýská	16	16
2119	blýskl	15	15
2120	bmi	13	14
2121	bmw	13	14
2122	bob	11	11
2123	boba	11	11
2124	bobbie	11	11
2125	bobby	13	11
2126	bobbyho	13	11
2127	bobe	11	11
2128	bobek	11	11
2129	bobem	13	14
2130	bobkový	15	15
2131	bobku	11	11
2132	bobošíková	18	18
2133	bobovi	11	11
2134	bobule	11	11
2135	boby	13	11
2136	bocích	17	17
2137	boční	17	17
2138	bočního	17	17
2139	bočních	17	17
2140	bočním	17	17
2141	bočními	17	17
2142	bod	11	11
2143	bodě	18	18
2144	bodech	12	12
2145	bodem	13	14
2146	bodl	11	11
2147	bodla	11	11
2148	bodnutí	17	17
2149	bodoval	11	11
2150	bodovali	11	11
2151	bodování	17	17
2152	bodovat	11	11
2153	bodové	18	18
2154	bodově	18	18
2155	bodového	18	18
2156	bodový	15	15
2157	bodových	15	15
2158	bodu	11	11
2159	bodů	11	11
2160	boduje	11	11
2161	bodům	13	14
2162	body	13	11
2163	boeing	12	12
2164	boha	11	11
2165	boháč	17	17
2166	boháče	17	17
2167	boháči	17	17
2168	boháčů	17	17
2169	bohatá	16	16
2170	bohaté	18	18
2171	bohatě	18	18
2172	bohatého	18	18
2173	bohatém	18	18
2174	bohatí	17	17
2175	bohatou	11	11
2176	bohatství	17	17
2177	bohatstvím	17	17
2178	bohatší	18	18
2179	bohatších	18	18
2180	bohatý	15	15
2181	bohatých	15	15
2182	bohatým	15	15
2183	bohatými	15	15
2184	bohdalová	16	16
2185	bohdan	12	12
2186	bohdana	12	12
2187	bohem	13	14
2188	bohemia	13	14
2189	bohemians	13	14
2190	bohnicích	17	17
2191	bohoslužba	15	15
2192	bohoslužbách	16	16
2193	bohoslužbě	18	18
2194	bohoslužbu	15	15
2195	bohoslužby	15	15
2196	bohoslužeb	15	15
2197	bohouš	18	18
2198	bohové	18	18
2199	bohu	11	11
2200	bohů	11	11
2201	bohům	13	14
2202	bohumil	13	14
2203	bohumila	13	14
2204	bohumín	17	17
2205	bohumína	17	17
2206	bohumíně	18	18
2207	bohumír	17	17
2208	bohuslav	11	11
2209	bohuslava	11	11
2210	bohužel	15	15
2211	bohy	13	11
2212	bohyně	18	18
2213	bohyni	13	12
2214	bohyní	17	17
2215	bochník	17	17
2216	boj	11	11
2217	boje	11	11
2218	bojem	13	14
2219	boji	11	11
2220	bojí	17	17
2221	bojích	17	17
2222	bojím	17	17
2223	bojíme	17	17
2224	bojíš	18	18
2225	bojiště	18	18
2226	bojišti	18	18
2227	bojíte	17	17
2228	bojkot	11	11
2229	bojová	16	16
2230	bojoval	11	11
2231	bojovala	11	11
2232	bojovali	11	11
2233	bojovalo	11	11
2234	bojovaly	13	11
2235	bojovat	11	11
2236	bojové	18	18
2237	bojového	18	18
2238	bojovné	18	18
2239	bojovně	18	18
2240	bojovnice	12	12
2241	bojovníci	17	17
2242	bojovník	17	17
2243	bojovníka	17	17
2244	bojovníků	17	17
2245	bojovníky	17	17
2246	bojovnost	12	12
2247	bojovný	15	15
2248	bojovou	11	11
2249	bojový	15	15
2250	bojových	15	15
2251	bojů	11	11
2252	bojuje	11	11
2253	bojujeme	13	14
2254	bojují	17	17
2255	bojující	17	17
2256	bojujících	17	17
2257	bok	11	11
2258	bokem	13	14
2259	boku	11	11
2260	boků	11	11
2261	boky	13	11
2262	bol	11	11
2263	bolavá	16	16
2264	bolavé	18	18
2265	bolek	11	11
2266	bolel	11	11
2267	bolela	11	11
2268	bolelo	11	11
2269	bolely	13	11
2270	boleslav	11	11
2271	boleslava	11	11
2272	boleslavi	11	11
2273	boleslaví	17	17
2274	bolest	11	11
2275	bolestech	12	12
2276	bolestem	13	14
2277	bolesti	11	11
2278	bolestí	17	17
2279	bolestivá	16	16
2280	bolestivé	18	18
2281	bolestivě	18	18
2282	bolestivý	15	15
2283	bolestivých	15	15
2284	bolestmi	13	14
2285	bolestná	16	16
2286	bolestné	18	18
2287	bolestně	18	18
2288	bolestnou	12	12
2289	bolestný	15	15
2290	bolestným	15	15
2291	bolet	11	11
2292	bolí	17	17
2293	bolívie	17	17
2294	bolívii	17	17
2295	bolševici	18	18
2296	bolt	11	11
2297	bomb	13	14
2298	bomba	13	14
2299	bombardéry	18	18
2300	bombardování	17	17
2301	bombou	13	14
2302	bombu	13	14
2303	bomby	13	14
2304	bon	12	12
2305	bonbon	12	12
2306	bonbonů	12	12
2307	bonbony	13	12
2308	bonbóny	19	19
2309	bond	12	12
2310	bonda	12	12
2311	bondy	13	12
2312	bonnie	12	12
2313	bonnu	12	12
2314	bonus	12	12
2315	bonusem	13	14
2316	bonusů	12	12
2317	bonusy	13	12
2318	book	11	11
2319	books	11	11
2320	boom	13	14
2321	boomu	13	14
2322	bor	11	11
2323	borce	12	12
2324	borci	12	12
2325	borců	12	12
2326	bordeaux	13	14
2327	bordel	11	11
2328	bordelu	11	11
2329	borec	12	12
2330	boris	11	11
2331	borise	11	11
2332	born	12	12
2333	borovic	12	12
2334	borovice	12	12
2335	boru	11	11
2336	borůvek	11	11
2337	borůvky	13	11
2338	borys	13	11
2339	bořek	16	16
2340	boří	17	17
2341	bos	11	11
2342	bosa	11	11
2343	bosá	16	16
2344	bosé	18	18
2345	bosch	12	12
2346	boskovic	12	12
2347	boskovice	12	12
2348	boskovicích	17	17
2349	bosně	18	18
2350	bosny	13	12
2351	boson	12	12
2352	bosonu	12	12
2353	bosony	13	12
2354	boss	11	11
2355	bosse	11	11
2356	boston	12	12
2357	bostonu	12	12
2358	bosý	15	15
2359	bosých	15	15
2360	bosýma	15	15
2361	bošňák	19	19
2362	bot	11	11
2363	bota	11	11
2364	botách	16	16
2365	botám	16	16
2366	botami	13	14
2367	botanická	16	16
2368	botanické	18	18
2369	botanických	15	15
2370	botanik	12	12
2371	botaniky	13	12
2372	botě	18	18
2373	botičky	17	17
2374	botky	13	11
2375	botou	11	11
2376	botu	11	11
2377	boty	13	11
2378	bouda	11	11
2379	boudě	18	18
2380	boudu	11	11
2381	boudy	13	11
2382	bouchačku	17	17
2383	bouchání	17	17
2384	bouchl	12	12
2385	bouchla	12	12
2386	bouchne	12	12
2387	boule	11	11
2388	bourání	17	17
2389	bourat	11	11
2390	bouřce	16	16
2391	bouře	16	16
2392	bouři	16	16
2393	bouří	17	17
2394	bouřka	16	16
2395	bouřky	16	16
2396	bouřlivá	16	16
2397	bouřlivé	18	18
2398	bouřlivě	18	18
2399	bouřlivý	16	16
2400	bouřlivých	16	16
2401	box	13	14
2402	boxech	13	14
2403	boxer	13	14
2404	boxu	13	14
2405	boxů	13	14
2406	boxy	13	14
2407	boy	13	11
2408	boyle	13	11
2409	boys	13	11
2410	bozi	11	14
2411	bože	15	15
2412	božena	15	15
2413	boženka	15	15
2414	boženy	15	15
2415	boží	17	17
2416	božího	17	17
2417	božích	17	17
2418	božím	17	17
2419	božská	16	16
2420	božské	18	18
2421	božského	18	18
2422	božskou	15	15
2423	božský	15	15
2424	božstev	15	15
2425	božstva	15	15
2426	božství	17	17
2427	božstvo	15	15
2428	bpm	13	14
2429	brabce	12	12
2430	brabec	12	12
2431	brad	11	11
2432	brada	11	11
2433	bradáčová	17	17
2434	bradavek	11	11
2435	bradavku	11	11
2436	bradavky	13	11
2437	bradě	18	18
2438	bradkou	11	11
2439	bradku	11	11
2440	bradley	13	11
2441	bradou	11	11
2442	bradu	11	11
2443	brady	13	11
2444	brácha	16	16
2445	brácho	16	16
2446	bráchou	16	16
2447	bráchovi	16	16
2448	bráchu	16	16
2449	bráchy	16	16
2450	bral	11	11
2451	brala	11	11
2452	brali	11	11
2453	bralo	11	11
2454	braly	13	11
2455	brambor	13	14
2456	bramborami	13	14
2457	bramborová	16	16
2458	bramborové	18	18
2459	bramborovou	13	14
2460	bramborový	15	15
2461	bramboru	13	14
2462	brambory	13	14
2463	brambůrky	13	14
2464	bran	12	12
2465	brán	16	16
2466	brána	16	16
2467	branami	13	14
2468	brance	12	12
2469	brand	12	12
2470	brandt	12	12
2471	brandy	13	12
2472	bráně	18	18
2473	branek	12	12
2474	bránění	18	18
2475	braní	17	17
2476	brání	17	17
2477	bránící	17	17
2478	bránil	16	16
2479	bránila	16	16
2480	bránili	16	16
2481	bránilo	16	16
2482	bránily	16	16
2483	bránit	16	16
2484	branka	12	12
2485	brankami	13	14
2486	brankář	16	16
2487	brankáře	16	16
2488	brankářem	16	16
2489	brankáři	16	16
2490	brankářů	16	16
2491	brankou	12	12
2492	brankové	18	18
2493	brankoviště	18	18
2494	brankovou	12	12
2495	branku	12	12
2496	branky	13	12
2497	bráno	16	16
2498	branou	12	12
2499	bránou	16	16
2500	bránu	16	16
2501	brány	16	16
2502	branže	15	15
2503	branži	15	15
2504	bráška	18	18
2505	brašně	18	18
2506	brašnu	18	18
2507	brašny	18	18
2508	brát	16	16
2509	bratislava	11	11
2510	bratislavě	18	18
2511	bratislavské	18	18
2512	bratislavy	13	11
2513	bratr	11	11
2514	bratra	11	11
2515	bratrance	12	12
2516	bratrancem	13	14
2517	bratranci	12	12
2518	bratranec	12	12
2519	bratrem	13	14
2520	bratrovi	11	11
2521	bratrské	18	18
2522	bratrstva	11	11
2523	bratrství	17	17
2524	bratrstvo	11	11
2525	bratru	11	11
2526	bratrů	11	11
2527	bratrům	13	14
2528	bratry	13	11
2529	bratře	16	16
2530	bratři	16	16
2531	bratří	17	17
2532	bratříčka	17	17
2533	braun	12	12
2534	brauner	12	12
2535	brával	16	16
2536	bravo	11	11
2537	bravurně	18	18
2538	brázda	16	16
2539	brázdí	17	17
2540	brázdy	16	16
2541	brazilci	12	14
2542	brazilec	12	14
2543	brazílie	17	17
2544	brazílii	17	17
2545	brazílií	17	17
2546	brazilské	18	18
2547	brazilského	18	18
2548	brazilském	18	18
2549	brazilský	15	15
2550	brečel	17	17
2551	brečela	17	17
2552	brečet	17	17
2553	brečí	17	17
2554	brečím	17	17
2555	brejku	11	11
2556	brejků	11	11
2557	brejky	13	11
2558	brejle	11	11
2559	breku	11	11
2560	brigáda	16	16
2561	brigádě	18	18
2562	brigádu	16	16
2563	brigády	16	16
2564	brilantně	18	18
2565	brilantní	17	17
2566	brit	11	11
2567	británie	16	16
2568	británii	16	16
2569	británií	17	17
2570	british	11	11
2571	britney	13	12
2572	britové	18	18
2573	britská	16	16
2574	britské	18	18
2575	britského	18	18
2576	britském	18	18
2577	britskému	18	18
2578	britskou	11	11
2579	britský	15	15
2580	britských	15	15
2581	britským	15	15
2582	britskými	15	15
2583	britští	18	18
2584	britů	11	11
2585	britům	13	14
2586	brity	13	11
2587	brna	12	12
2588	brňané	19	19
2589	brně	18	18
2590	brnem	13	14
2591	brnění	18	18
2592	brněnská	18	18
2593	brněnské	18	18
2594	brněnského	18	18
2595	brněnském	18	18
2596	brněnskou	18	18
2597	brněnsku	18	18
2598	brněnský	18	18
2599	brněnských	18	18
2600	brněnským	18	18
2601	brněnští	18	18
2602	brno	12	12
2603	brnu	12	12
2604	brod	11	11
2605	brodě	18	18
2606	brodit	11	11
2607	brodu	11	11
2608	brokolice	12	12
2609	brokovnici	12	12
2610	bronz	12	14
2611	bronzová	16	16
2612	bronzové	18	18
2613	bronzovou	12	14
2614	bronzový	15	15
2615	bronzových	15	15
2616	bronzu	12	14
2617	brooklynu	13	12
2618	broskve	11	11
2619	brouci	12	12
2620	broučku	17	17
2621	brouk	11	11
2622	brouka	11	11
2623	brouků	11	11
2624	brouky	13	11
2625	brousí	17	17
2626	brousit	11	11
2627	broušené	18	18
2628	broušení	18	18
2629	brož	15	15
2630	brožová	16	16
2631	brožury	15	15
2632	bručel	17	17
2633	bruno	12	12
2634	bruntál	16	16
2635	bruntálu	16	16
2636	brusel	11	11
2637	bruselu	11	11
2638	brusle	11	11
2639	bruslení	17	17
2640	bruslí	17	17
2641	bruslích	17	17
2642	bruslit	11	11
2643	brutálně	18	18
2644	brutální	17	17
2645	brutálním	17	17
2646	brýle	15	15
2647	brýlemi	15	15
2648	brýlí	17	17
2649	brýlích	17	17
2650	brzd	11	14
2651	brzda	11	14
2652	brzdění	18	18
2653	brzdí	17	17
2654	brzdit	11	14
2655	brzdou	11	14
2656	brzdové	18	18
2657	brzdu	11	14
2658	brzdy	13	14
2659	brzké	18	18
2660	brzký	15	15
2661	brzo	11	14
2662	brzobohatý	15	15
2663	brzy	13	14
2664	břeclav	16	16
2665	břeclavi	16	16
2666	břeclavské	18	18
2667	břeclavsko	16	16
2668	břeclavsku	16	16
2669	břečťan	19	19
2670	břeh	16	16
2671	břehem	16	16
2672	břehu	16	16
2673	břehů	16	16
2674	břehům	16	16
2675	břehy	16	16
2676	břemen	16	16
2677	břemena	16	16
2678	břemene	16	16
2679	břemenem	16	16
2680	břemeno	16	16
2681	břetislav	16	16
2682	břevno	16	16
2683	březen	16	16
2684	březích	17	17
2685	březina	16	16
2686	březiny	16	16
2687	března	16	16
2688	březnové	18	18
2689	březnovém	18	18
2690	březnu	16	16
2691	březosti	16	16
2692	březová	16	16
2693	březové	18	18
2694	břidlic	16	16
2695	břidlice	16	16
2696	břicha	16	16
2697	břichem	16	16
2698	břicho	16	16
2699	břichu	16	16
2700	břímě	18	18
2701	břiše	18	18
2702	bříška	18	18
2703	bříškem	18	18
2704	bříško	18	18
2705	bříšku	18	18
2706	břišní	18	18
2707	břitva	16	16
2708	bříza	17	17
2709	břízy	17	17
2710	buben	12	12
2711	bubeník	17	17
2712	bubínek	17	17
2713	bubínky	17	17
2714	bublá	16	16
2715	bublin	12	12
2716	bublina	12	12
2717	bublině	18	18
2718	bublinek	12	12
2719	bublinky	13	12
2720	bublinu	12	12
2721	bubliny	13	12
2722	bubnoval	12	12
2723	bubnování	17	17
2724	bubnu	12	12
2725	bubnů	12	12
2726	bubny	13	12
2727	buď	19	19
2728	budapešť	19	19
2729	budapešti	18	18
2730	budce	12	12
2731	buddha	11	11
2732	buddhismu	13	14
2733	buddhismus	13	14
2734	buddhistické	18	18
2735	buddhy	13	11
2736	bude	11	11
2737	budějovic	18	18
2738	budějovice	18	18
2739	budějovicemi	18	18
2740	budějovicích	18	18
2741	budem	13	14
2742	budeme	13	14
2743	budeš	18	18
2744	budete	11	11
2745	budí	17	17
2746	budík	17	17
2747	budíku	17	17
2748	budil	11	11
2749	budila	11	11
2750	budily	13	11
2751	budit	11	11
2752	budiž	15	15
2753	budižkničemu	17	17
2754	budky	13	11
2755	buďme	19	19
2756	budou	11	11
2757	budoucí	17	17
2758	budoucího	17	17
2759	budoucích	17	17
2760	budoucím	17	17
2761	budoucímu	17	17
2762	budoucna	12	12
2763	budoucnost	12	12
2764	budoucnosti	12	12
2765	budoucností	17	17
2766	budoucnu	12	12
2767	budov	11	11
2768	budova	11	11
2769	budovách	16	16
2770	budoval	11	11
2771	budovala	11	11
2772	budovali	11	11
2773	budovami	13	14
2774	budována	16	16
2775	budované	18	18
2776	budování	17	17
2777	budováním	17	17
2778	budovány	16	16
2779	budovat	11	11
2780	budově	18	18
2781	budovou	11	11
2782	budovu	11	11
2783	budovy	13	11
2784	buďte	19	19
2785	buďto	19	19
2786	budu	11	11
2787	buduje	11	11
2788	budují	17	17
2789	budvar	11	11
2790	buenos	12	12
2791	bufet	11	11
2792	bufetu	11	11
2793	buffalo	11	11
2794	bůh	11	11
2795	bůhví	17	17
2796	bůhvíco	17	17
2797	bůhvíjak	17	17
2798	bůhvíproč	17	17
2799	buchtička	17	17
2800	buchtu	12	12
2801	buchty	13	12
2802	building	12	12
2803	bujení	17	17
2804	bují	17	17
2805	bujné	18	18
2806	bujnou	12	12
2807	buk	11	11
2808	buku	11	11
2809	buky	13	11
2810	bulharska	11	11
2811	bulharské	18	18
2812	bulharsko	11	11
2813	bulharsku	11	11
2814	bull	11	11
2815	bullu	11	11
2816	bully	13	11
2817	bulovce	12	12
2818	bulvár	16	16
2819	bulvární	17	17
2820	bulvárních	17	17
2821	bulváru	16	16
2822	bulvy	13	11
2823	buly	13	11
2824	bum	13	14
2825	buňce	19	19
2826	bunda	12	12
2827	bundách	16	16
2828	bundě	18	18
2829	bundu	12	12
2830	bundy	13	12
2831	buněčné	18	18
2832	buněčného	18	18
2833	buněčnou	18	18
2834	buněčných	18	18
2835	buněk	18	18
2836	bungalovu	12	12
2837	buňka	19	19
2838	buňkách	19	19
2839	buňkám	19	19
2840	buňkami	19	19
2841	bunkr	12	12
2842	bunkru	12	12
2843	buňku	19	19
2844	buňky	19	19
2845	burácení	17	17
2846	burčák	17	17
2847	burda	11	11
2848	bureš	18	18
2849	burger	11	11
2850	burian	12	12
2851	buriana	12	12
2852	burke	11	11
2853	burkhalter	11	11
2854	burroughs	11	11
2855	bursík	17	17
2856	burton	12	12
2857	burundi	12	12
2858	burza	11	14
2859	burzách	16	16
2860	burze	11	14
2861	burzovní	17	17
2862	burzu	11	14
2863	burzy	13	14
2864	buržoazie	15	15
2865	buržoazní	17	17
2866	bush	11	11
2867	bushe	11	11
2868	business	12	12
2869	bušení	18	18
2870	buší	18	18
2871	bušícím	18	18
2872	bušil	18	18
2873	bušila	18	18
2874	bušilo	18	18
2875	bušit	18	18
2876	butch	12	12
2877	butik	11	11
2878	butiku	11	11
2879	buzení	17	17
2880	bwv	11	11
2881	býci	15	15
2882	býčí	17	17
2883	bydlel	13	11
2884	bydlela	13	11
2885	bydleli	13	11
2886	bydlely	13	11
2887	bydlení	17	17
2888	bydlením	17	17
2889	bydlet	13	11
2890	bydlí	17	17
2891	bydlící	17	17
2892	bydlím	17	17
2893	bydlíme	17	17
2894	bydlíš	18	18
2895	bydliště	18	18
2896	bydlišti	18	18
2897	bydlíte	17	17
2898	bych	13	12
2899	bychom	13	14
2900	býk	15	15
2901	býka	15	15
2902	byl	13	11
2903	byla	13	11
2904	byli	13	11
2905	bylin	13	12
2906	bylina	13	12
2907	bylinách	16	16
2908	bylinek	13	12
2909	bylinkami	13	14
2910	bylinkové	18	18
2911	bylinky	13	12
2912	bylinné	18	18
2913	byliny	13	12
2914	bylo	13	11
2915	byly	13	11
2916	byrokracie	13	12
2917	byrokracii	13	12
2918	byrokratické	18	18
2919	bys	13	11
2920	bysme	13	14
2921	byste	13	11
2922	bystrá	16	16
2923	bystré	18	18
2924	bystrý	15	15
2925	bystře	16	16
2926	bystřice	16	16
2927	bystřici	16	16
2928	byt	13	11
2929	byť	19	19
2930	být	15	15
2931	bytě	18	18
2932	bytech	13	12
2933	bytem	13	14
2934	bytí	17	17
2935	býti	15	15
2936	bytím	17	17
2937	bytost	13	11
2938	bytostem	13	14
2939	bytosti	13	11
2940	bytostí	17	17
2941	bytostmi	13	14
2942	bytostně	18	18
2943	bytová	16	16
2944	bytové	18	18
2945	bytového	18	18
2946	bytovém	18	18
2947	bytovou	13	11
2948	bytový	15	15
2949	bytových	15	15
2950	bytu	13	11
2951	bytů	13	11
2952	byty	13	11
2953	bývá	16	16
2954	bývají	17	17
2955	býval	15	15
2956	bývala	15	15
2957	bývalá	16	16
2958	bývalé	18	18
2959	bývalého	18	18
2960	bývalém	18	18
2961	bývalému	18	18
2962	bývali	15	15
2963	bývalí	17	17
2964	bývalo	15	15
2965	bývalou	15	15
2966	bývaly	15	15
2967	bývalý	15	15
2968	bývalých	15	15
2969	bývalým	15	15
2970	bývalými	15	15
2971	bývám	16	16
2972	byznys	13	14
2973	byznysem	13	14
2974	byznysu	13	14
2975	bzučení	17	17
2976	bzučí	17	17
2977	caesar	12	12
2978	café	18	18
2979	california	12	12
2980	cambridge	13	14
2981	cár	16	16
2982	cáry	16	16
2983	cédéčka	18	18
2984	cédéčko	18	18
2985	cederna	12	12
2986	cedule	12	12
2987	ceduli	12	12
2988	cedulí	17	17
2989	cedulka	12	12
2990	cedulku	12	12
2991	cejtil	12	12
2992	cela	12	12
2993	celá	16	16
2994	celách	16	16
2995	celé	18	18
2996	celebrit	12	12
2997	celebrita	12	12
2998	celebritou	12	12
2999	celebrity	13	12
3000	celého	18	18
3001	celej	12	12
3002	celek	12	12
3003	celém	18	18
3004	celému	18	18
3005	celer	12	12
3006	celeru	12	12
3007	celetné	18	18
3008	celí	17	17
3009	celia	12	12
3010	celistvé	18	18
3011	celistvost	12	12
3012	celistvosti	12	12
3013	celistvý	15	15
3014	celkem	13	14
3015	celková	16	16
3016	celkové	18	18
3017	celkově	18	18
3018	celkového	18	18
3019	celkovém	18	18
3020	celkovému	18	18
3021	celkovou	12	12
3022	celkový	15	15
3023	celkových	15	15
3024	celkovým	15	15
3025	celku	12	12
3026	celků	12	12
3027	celky	13	12
3028	celní	17	17
3029	celníci	17	17
3030	celního	17	17
3031	celních	17	17
3032	celníků	17	17
3033	celodenní	17	17
3034	celodenním	17	17
3035	celoevropské	18	18
3036	celonárodní	17	17
3037	celoročně	18	18
3038	celoroční	17	17
3039	celostátní	17	17
3040	celostátního	17	17
3041	celostátních	17	17
3042	celostátním	17	17
3043	celosvětová	18	18
3044	celosvětové	18	18
3045	celosvětově	18	18
3046	celosvětového	18	18
3047	celosvětovém	18	18
3048	celosvětovou	18	18
3049	celosvětový	18	18
3050	celou	12	12
3051	celovečerní	17	17
3052	celozrnné	18	18
3053	celoživotní	17	17
3054	celoživotního	17	17
3055	celoživotním	17	17
3056	celsia	12	12
3057	celu	12	12
3058	cely	13	12
3059	celý	15	15
3060	celýho	15	15
3061	celých	15	15
3062	celým	15	15
3063	celými	15	15
3064	cement	13	14
3065	cementu	13	14
3066	cen	12	12
3067	cena	12	12
3068	cenách	16	16
3069	cenám	16	16
3070	cenami	13	14
3071	ceně	18	18
3072	ceněné	18	18
3073	cení	17	17
3074	ceník	17	17
3075	ceníku	17	17
3076	cenil	12	12
3077	cením	17	17
3078	cenná	16	16
3079	cenné	18	18
3080	cenného	18	18
3081	cennější	18	18
3082	cennosti	12	12
3083	cennou	12	12
3084	cenný	15	15
3085	cenných	15	15
3086	cenným	15	15
3087	cennými	15	15
3088	cenou	12	12
3089	cenová	16	16
3090	cenové	18	18
3091	cenově	18	18
3092	cenovou	12	12
3093	cenový	15	15
3094	cenových	15	15
3095	cent	12	12
3096	centauri	12	12
3097	center	12	12
3098	centimetr	13	14
3099	centimetru	13	14
3100	centimetrů	13	14
3101	centimetry	13	14
3102	centr	12	12
3103	centra	12	12
3104	central	12	12
3105	centrála	16	16
3106	centrále	16	16
3107	centralizace	12	14
3108	centrálně	18	18
3109	centrální	17	17
3110	centrálního	17	17
3111	centrálních	17	17
3112	centrálním	17	17
3113	centrálu	16	16
3114	centrály	16	16
3115	centre	12	12
3116	centrech	12	12
3117	centrem	13	14
3118	centru	12	12
3119	centrum	13	14
3120	centry	13	12
3121	centů	12	12
3122	century	13	12
3123	cenu	12	12
3124	ceny	13	12
3125	cenzura	12	14
3126	cenzury	13	14
3127	ceremoniál	16	16
3128	ceremoniálu	16	16
3129	ceremonie	13	14
3130	ceresit	12	12
3131	cerevisiae	12	12
3132	cern	12	12
3133	cernu	12	12
3134	certifikace	12	12
3135	certifikaci	12	12
3136	certifikací	17	17
3137	certifikační	17	17
3138	certifikát	16	16
3139	certifikátu	16	16
3140	certifikátů	16	16
3141	certifikáty	16	16
3142	certifikované	18	18
3143	certifikovaných	15	15
3144	cest	12	12
3145	cesta	12	12
3146	cestách	16	16
3147	cestám	16	16
3148	cestami	13	14
3149	cestě	18	18
3150	cestičce	17	17
3151	cestička	17	17
3152	cestičku	17	17
3153	cestičky	17	17
3154	cestou	12	12
3155	cestoval	12	12
3156	cestovala	12	12
3157	cestovali	12	12
3158	cestování	17	17
3159	cestováním	17	17
3160	cestovat	12	12
3161	cestovatel	12	12
3162	cestovatele	12	12
3163	cestovatelé	18	18
3164	cestovatelů	12	12
3165	cestovky	13	12
3166	cestovní	17	17
3167	cestovního	17	17
3168	cestovních	17	17
3169	cestovním	17	17
3170	cestu	12	12
3171	cestuje	12	12
3172	cestujete	12	12
3173	cestují	17	17
3174	cestující	17	17
3175	cestujícího	17	17
3176	cestujících	17	17
3177	cestujícím	17	17
3178	cestujícími	17	17
3179	cesty	13	12
3180	cév	18	18
3181	cévní	18	18
3182	cévních	18	18
3183	cévy	18	18
3184	cia	12	12
3185	cibule	12	12
3186	cibuli	12	12
3187	cibulí	17	17
3188	cibulka	12	12
3189	cibulkou	12	12
3190	cibulku	12	12
3191	cibulky	13	12
3192	cicely	13	12
3193	cicero	12	12
3194	cidlinou	12	12
3195	ciferník	17	17
3196	cigára	16	16
3197	cigaret	12	12
3198	cigareta	12	12
3199	cigaretami	13	14
3200	cigaretou	12	12
3201	cigaretového	18	18
3202	cigaretový	15	15
3203	cigaretu	12	12
3204	cigarety	13	12
3205	cigáro	16	16
3206	cihel	12	12
3207	cihelné	18	18
3208	cihla	12	12
3209	cihlové	18	18
3210	cihlovou	12	12
3211	cihlu	12	12
3212	cihly	13	12
3213	cikánka	16	16
3214	cíl	17	17
3215	cíle	17	17
3216	cílech	17	17
3217	cílek	17	17
3218	cílem	17	17
3219	cílená	17	17
3220	cílené	18	18
3221	cíleně	18	18
3222	cílenou	17	17
3223	cílevědomě	18	18
3224	cíli	17	17
3225	cílí	17	17
3226	cílová	17	17
3227	cílové	18	18
3228	cílového	18	18
3229	cílovou	17	17
3230	cílový	17	17
3231	cílových	17	17
3232	cílů	17	17
3233	cílům	17	17
3234	cimbálová	16	16
3235	cimbura	13	14
3236	cimbuří	17	17
3237	cimrman	13	14
3238	cimrmana	13	14
3239	cindy	13	12
3240	cinema	13	14
3241	cinkání	17	17
3242	cinkot	12	12
3243	cínu	17	17
3244	cíp	17	17
3245	cípu	17	17
3246	cípy	17	17
3247	cir	12	12
3248	církev	17	17
3249	církevní	17	17
3250	církevního	17	17
3251	církevních	17	17
3252	církevním	17	17
3253	cirkulace	12	12
3254	cirkulaci	12	12
3255	cirkus	12	12
3256	cirkusu	12	12
3257	církve	17	17
3258	církvemi	17	17
3259	církvi	17	17
3260	církví	17	17
3261	církvím	17	17
3262	císař	17	17
3263	císaře	17	17
3264	císařem	17	17
3265	císaři	17	17
3266	císařské	18	18
3267	císařského	18	18
3268	císařském	18	18
3269	císařský	17	17
3270	císařských	17	17
3271	císařským	17	17
3272	císařství	17	17
3273	císařští	18	18
3274	císařů	17	17
3275	cisterny	13	12
3276	cit	12	12
3277	citace	12	12
3278	citací	17	17
3279	citát	16	16
3280	citátu	16	16
3281	citáty	16	16
3282	citech	12	12
3283	citelné	18	18
3284	citelně	18	18
3285	citem	13	14
3286	cítění	18	18
3287	cítí	17	17
3288	cítil	17	17
3289	cítila	17	17
3290	cítili	17	17
3291	cítily	17	17
3292	cítím	17	17
3293	cítíme	17	17
3294	cítíš	18	18
3295	cítit	17	17
3296	cítíte	17	17
3297	citlivá	16	16
3298	citlivé	18	18
3299	citlivě	18	18
3300	citlivého	18	18
3301	citlivější	18	18
3302	citliví	17	17
3303	citlivost	12	12
3304	citlivosti	12	12
3305	citlivostí	17	17
3306	citlivou	12	12
3307	citlivý	15	15
3308	citlivých	15	15
3309	citlivým	15	15
3310	citová	16	16
3311	citoval	12	12
3312	citovala	12	12
3313	citované	18	18
3314	citovaný	15	15
3315	citovaných	15	15
3316	citovat	12	12
3317	citové	18	18
3318	citově	18	18
3319	citového	18	18
3320	citovou	12	12
3321	citový	15	15
3322	citových	15	15
3323	citroën	-1	-1
3324	citron	12	12
3325	citronem	13	14
3326	citronová	16	16
3327	citronové	18	18
3328	citronovou	12	12
3329	citronu	12	12
3330	citu	12	12
3331	citů	12	12
3332	cituje	12	12
3333	cituji	12	12
3334	city	13	12
3335	cívek	17	17
3336	civěl	18	18
3337	civěla	18	18
3338	civěli	18	18
3339	civí	17	17
3340	civilisté	18	18
3341	civilistů	12	12
3342	civilisty	13	12
3343	civilizace	12	14
3344	civilizaci	12	14
3345	civilizací	17	17
3346	civilizační	17	17
3347	civilizačních	17	17
3348	civilizované	18	18
3349	civilní	17	17
3350	civilního	17	17
3351	civilních	17	17
3352	civilním	17	17
3353	civilu	12	12
3354	cívka	17	17
3355	cívkou	17	17
3356	cívku	17	17
3357	cívky	17	17
3358	cize	12	14
3359	cizí	17	17
3360	cizího	17	17
3361	cizích	17	17
3362	cizím	17	17
3363	cizími	17	17
3364	cizímu	17	17
3365	cizince	12	14
3366	cizincem	13	14
3367	cizinci	12	14
3368	cizinců	12	14
3369	cizincům	13	14
3370	cizině	18	18
3371	cizinec	12	14
3372	cizinecké	18	18
3373	cizinka	12	14
3374	ciziny	13	14
3375	clona	12	12
3376	clonou	12	12
3377	clonu	12	12
3378	clony	13	12
3379	co	12	12
3380	coby	13	12
3381	cokoli	12	12
3382	cokoliv	12	12
3383	colorado	12	12
3384	columbia	13	14
3385	comfort	13	14
3386	cop	12	12
3387	copak	12	12
3388	copyright	13	12
3389	cosi	12	12
3390	couvá	16	16
3391	couval	12	12
3392	couvání	17	17
3393	couvat	12	12
3394	couvl	12	12
3395	couvla	12	12
3396	couvnout	12	12
3397	což	15	15
3398	cože	15	15
3399	cožpak	15	15
3400	ctí	17	17
3401	ctihodný	15	15
3402	ctít	17	17
3403	ctižádost	16	16
3404	ctižádostivý	16	16
3405	ctnost	12	12
3406	ctnosti	12	12
3407	ctností	17	17
3408	cucky	13	12
3409	cukety	13	12
3410	cukr	12	12
3411	cukrárně	18	18
3412	cukrárny	16	16
3413	cukrem	13	14
3414	cukrové	18	18
3415	cukroví	17	17
3416	cukrovka	12	12
3417	cukrovkou	12	12
3418	cukrovku	12	12
3419	cukrovky	13	12
3420	cukrovou	12	12
3421	cukru	12	12
3422	cukrů	12	12
3423	cukry	13	12
3424	cvak	12	12
3425	cvakání	17	17
3426	cvaknutí	17	17
3427	cvičební	17	17
3428	cvičení	17	17
3429	cvičeních	17	17
3430	cvičením	17	17
3431	cvičí	17	17
3432	cvičil	17	17
3433	cvičila	17	17
3434	cvičili	17	17
3435	cvičit	17	17
3436	cvik	12	12
3437	cviku	12	12
3438	cviků	12	12
3439	cviky	13	12
3440	cvok	12	12
3441	cyklické	18	18
3442	cyklický	15	15
3443	cyklista	13	12
3444	cyklisté	18	18
3445	cyklistice	13	12
3446	cyklistické	18	18
3447	cyklistického	18	18
3448	cyklistickou	13	12
3449	cyklistický	15	15
3450	cyklistických	15	15
3451	cyklistika	13	12
3452	cyklistiku	13	12
3453	cyklistiky	13	12
3454	cyklistu	13	12
3455	cyklistů	13	12
3456	cyklistům	13	14
3457	cyklisty	13	12
3458	cyklostezek	13	14
3459	cyklostezka	13	14
3460	cyklostezku	13	14
3461	cyklostezky	13	14
3462	cyklus	13	12
3463	cynický	15	15
3464	cyril	13	12
3465	cyrila	13	12
3466	čaj	17	17
3467	čaje	17	17
3468	čajem	17	17
3469	čaji	17	17
3470	čajkovského	18	18
3471	čajové	18	18
3472	čajových	17	17
3473	čajů	17	17
3474	čalounění	18	18
3475	čáp	17	17
3476	čapek	17	17
3477	čapka	17	17
3478	čapku	17	17
3479	čar	17	17
3480	čára	17	17
3481	čarami	17	17
3482	čárku	17	17
3483	čárky	17	17
3484	čaroděj	18	18
3485	čaroděje	18	18
3486	čarodějnic	18	18
3487	čarodějnice	18	18
3488	čarodějnici	18	18
3489	čarodějnictví	18	18
3490	čarou	17	17
3491	čárou	17	17
3492	čárových	17	17
3493	čáru	17	17
3494	čáry	17	17
3495	čáře	17	17
3496	čas	17	17
3497	čase	17	17
3498	časech	17	17
3499	časem	17	17
3500	čáslav	17	17
3501	čáslavi	17	17
3502	čáslavská	17	17
3503	časné	18	18
3504	časně	18	18
3505	časného	18	18
3506	časném	18	18
3507	časných	17	17
3508	časopis	17	17
3509	časopise	17	17
3510	časopisech	17	17
3511	časopisem	17	17
3512	časopisu	17	17
3513	časopisů	17	17
3514	časopisy	17	17
3515	časoprostoru	17	17
3516	časová	17	17
3517	časování	17	17
3518	časové	18	18
3519	časově	18	18
3520	časového	18	18
3521	časovém	18	18
3522	časovému	18	18
3523	časovou	17	17
3524	časový	17	17
3525	časových	17	17
3526	časovým	17	17
3527	časovými	17	17
3528	část	17	17
3529	častá	17	17
3530	částce	17	17
3531	časté	18	18
3532	částeček	17	17
3533	částečky	17	17
3534	částečná	17	17
3535	částečné	18	18
3536	částečně	18	18
3537	částečného	18	18
3538	částečnou	17	17
3539	částečný	17	17
3540	částečným	17	17
3541	častého	18	18
3542	částech	17	17
3543	častěji	18	18
3544	častější	18	18
3545	častějším	18	18
3546	částek	17	17
3547	částem	17	17
3548	části	17	17
3549	částí	17	17
3550	částic	17	17
3551	částice	17	17
3552	částicemi	17	17
3553	částici	17	17
3554	částicové	18	18
3555	částka	17	17
3556	částkou	17	17
3557	částku	17	17
3558	částky	17	17
3559	částmi	17	17
3560	často	17	17
3561	častokrát	17	17
3562	častou	17	17
3563	častý	17	17
3564	častých	17	17
3565	častým	17	17
3566	častými	17	17
3567	času	17	17
3568	časů	17	17
3569	časy	17	17
3570	čau	17	17
3571	čavu	17	17
3572	čediče	17	17
3573	čedičové	18	18
3574	čedičových	17	17
3575	čeho	17	17
3576	čehokoli	17	17
3577	čehokoliv	17	17
3578	čehosi	17	17
3579	čehož	17	17
3580	čech	17	17
3581	čecha	17	17
3582	čechách	17	17
3583	čechem	17	17
3584	čechová	17	17
3585	čechovi	17	17
3586	čechů	17	17
3587	čechům	17	17
3588	čechy	17	17
3589	čeká	17	17
3590	čekací	17	17
3591	čekají	17	17
3592	čekající	17	17
3593	čekajících	17	17
3594	čekal	17	17
3595	čekala	17	17
3596	čekali	17	17
3597	čekalo	17	17
3598	čekaly	17	17
3599	čekám	17	17
3600	čekáme	17	17
3601	čekání	17	17
3602	čekáním	17	17
3603	čekárně	18	18
3604	čekárny	17	17
3605	čekáš	18	18
3606	čekat	17	17
3607	čekáte	17	17
3608	čela	17	17
3609	čelákovice	17	17
3610	čele	17	17
3611	čeleď	19	19
3612	čeledi	17	17
3613	čelem	17	17
3614	čelenku	17	17
3615	čelí	17	17
3616	čelil	17	17
3617	čelila	17	17
3618	čelili	17	17
3619	čelist	17	17
3620	čelisti	17	17
3621	čelistí	17	17
3622	čelistmi	17	17
3623	čelit	17	17
3624	čelně	18	18
3625	čelní	17	17
3626	čelního	17	17
3627	čelních	17	17
3628	čelním	17	17
3629	čelo	17	17
3630	čelu	17	17
3631	čem	17	17
3632	čemkoli	17	17
3633	čemkoliv	17	17
3634	čemsi	17	17
3635	čemu	17	17
3636	čemukoli	17	17
3637	čemuž	17	17
3638	čemž	17	17
3639	čeněk	18	18
3640	čep	17	17
3641	čepec	17	17
3642	čepel	17	17
3643	čepele	17	17
3644	čepelí	17	17
3645	čepice	17	17
3646	čepici	17	17
3647	čepicí	17	17
3648	čeps	17	17
3649	čermák	17	17
3650	čermáka	17	17
3651	čerň	19	19
3652	černá	17	17
3653	černé	18	18
3654	černě	18	18
3655	černého	18	18
3656	černej	17	17
3657	černém	18	18
3658	černému	18	18
3659	černí	17	17
3660	černo	17	17
3661	černobílá	17	17
3662	černobíle	17	17
3663	černobílé	18	18
3664	černobílého	18	18
3665	černobílou	17	17
3666	černobílý	17	17
3667	černobílých	17	17
3668	černobylu	17	17
3669	černoch	17	17
3670	černocha	17	17
3671	černochů	17	17
3672	černochy	17	17
3673	černoši	18	18
3674	černošské	18	18
3675	černou	17	17
3676	černý	17	17
3677	černých	17	17
3678	černým	17	17
3679	černýma	17	17
3680	černými	17	17
3681	čerpá	17	17
3682	čerpací	17	17
3683	čerpacích	17	17
3684	čerpadel	17	17
3685	čerpadla	17	17
3686	čerpadlem	17	17
3687	čerpadlo	17	17
3688	čerpají	17	17
3689	čerpal	17	17
3690	čerpala	17	17
3691	čerpali	17	17
3692	čerpání	17	17
3693	čerpáním	17	17
3694	čerpat	17	17
3695	čerstvá	17	17
3696	čerstvé	18	18
3697	čerstvě	18	18
3698	čerstvého	18	18
3699	čerstvém	18	18
3700	čerstvou	17	17
3701	čerstvý	17	17
3702	čerstvých	17	17
3703	čerstvým	17	17
3704	čerstvými	17	17
3705	čert	17	17
3706	čerta	17	17
3707	čerti	17	17
3708	čertova	17	17
3709	čertu	17	17
3710	čertů	17	17
3711	červ	17	17
3712	červa	17	17
3713	červen	17	17
3714	červeň	19	19
3715	červená	17	17
3716	července	17	17
3717	červenci	17	17
3718	červencové	18	18
3719	červené	18	18
3720	červeně	18	18
3721	červenec	17	17
3722	červeného	18	18
3723	červeném	18	18
3724	červení	17	17
3725	červenka	17	17
3726	červenohnědé	18	18
3727	červenou	17	17
3728	červený	17	17
3729	červených	17	17
3730	červeným	17	17
3731	červenými	17	17
3732	červi	17	17
3733	června	17	17
3734	červnové	18	18
3735	červnovém	18	18
3736	červnu	17	17
3737	červů	17	17
3738	červy	17	17
3739	česka	17	17
3740	česká	17	17
3741	české	18	18
3742	českého	18	18
3743	českem	17	17
3744	českém	18	18
3745	českému	18	18
3746	česko	17	17
3747	českobudějovické	18	18
3748	českobudějovického	18	18
3749	českolipsku	17	17
3750	českomoravské	18	18
3751	českomoravský	17	17
3752	československa	17	17
3753	československá	17	17
3754	československé	18	18
3755	československého	18	18
3756	československo	17	17
3757	československou	17	17
3758	československu	17	17
3759	československý	17	17
3760	československých	17	17
3761	československým	17	17
3762	českoslovenští	18	18
3763	českou	17	17
3764	česku	17	17
3765	česky	17	17
3766	český	17	17
3767	českých	17	17
3768	českým	17	17
3769	českými	17	17
3770	česnek	17	17
3771	česnekem	17	17
3772	česneku	17	17
3773	čest	17	17
3774	čestmír	17	17
3775	čestná	17	17
3776	čestné	18	18
3777	čestně	18	18
3778	čestného	18	18
3779	čestnou	17	17
3780	čestný	17	17
3781	čestných	17	17
3782	čestným	17	17
3783	češi	18	18
3784	češích	18	18
3785	češka	18	18
3786	češky	18	18
3787	čeští	18	18
3788	čeština	18	18
3789	češtině	18	18
3790	češtinou	18	18
3791	češtinu	18	18
3792	češtiny	18	18
3793	čet	17	17
3794	četa	17	17
3795	četba	17	17
3796	četbě	18	18
3797	četbou	17	17
3798	četbu	17	17
3799	četby	17	17
3800	četě	18	18
3801	četl	17	17
3802	četla	17	17
3803	četli	17	17
3804	četly	17	17
3805	četná	17	17
3806	četné	18	18
3807	četníci	17	17
3808	četník	17	17
3809	četnost	17	17
3810	četnosti	17	17
3811	četností	17	17
3812	četných	17	17
3813	četným	17	17
3814	četnými	17	17
3815	četu	17	17
3816	čety	17	17
3817	čez	17	17
3818	čfl	17	17
3819	čching	17	17
3820	čidel	17	17
3821	čidla	17	17
3822	čidlo	17	17
3823	číhá	17	17
3824	číhají	17	17
3825	číhal	17	17
3826	číhat	17	17
3827	čich	17	17
3828	čichu	17	17
3829	čijo	17	17
3830	čile	17	17
3831	čili	17	17
3832	čilý	17	17
3833	čím	17	17
3834	čímkoli	17	17
3835	čímsi	17	17
3836	čímž	17	17
3837	čin	17	17
3838	čína	17	17
3839	číňan	19	19
3840	číňané	19	19
3841	číňanů	19	19
3842	číňany	19	19
3843	číně	18	18
3844	činech	17	17
3845	činem	17	17
3846	činění	18	18
3847	činí	17	17
3848	činidla	17	17
3849	činidlo	17	17
3850	činil	17	17
3851	činila	17	17
3852	činili	17	17
3853	činilo	17	17
3854	činily	17	17
3855	činit	17	17
3856	činitel	17	17
3857	činitele	17	17
3858	činitelé	18	18
3859	činitelem	17	17
3860	činiteli	17	17
3861	činitelů	17	17
3862	činky	17	17
3863	činné	18	18
3864	činnost	17	17
3865	činnostech	17	17
3866	činnostem	17	17
3867	činnosti	17	17
3868	činností	17	17
3869	činnostmi	17	17
3870	činný	17	17
3871	činných	17	17
3872	činohry	17	17
3873	čínou	17	17
3874	čínská	17	17
3875	čínské	18	18
3876	čínského	18	18
3877	čínském	18	18
3878	čínskou	17	17
3879	čínsky	17	17
3880	čínský	17	17
3881	čínských	17	17
3882	čínským	17	17
3883	čínskými	17	17
3884	čínští	18	18
3885	činu	17	17
3886	činů	17	17
3887	čínu	17	17
3888	činům	17	17
3889	činy	17	17
3890	číny	17	17
3891	činžáku	17	17
3892	činžovní	17	17
3893	čip	17	17
3894	čipem	17	17
3895	čípku	17	17
3896	čipu	17	17
3897	čipů	17	17
3898	čipy	17	17
3899	čirá	17	17
3900	čiré	18	18
3901	čirého	18	18
3902	čirou	17	17
3903	čirý	17	17
3904	čísel	17	17
3905	číselné	18	18
3906	číselných	17	17
3907	čísi	17	17
3908	čísla	17	17
3909	čísle	17	17
3910	číslech	17	17
3911	číslem	17	17
3912	číslic	17	17
3913	číslice	17	17
3914	číslicemi	17	17
3915	číslo	17	17
3916	číslu	17	17
3917	číslům	17	17
3918	čísly	17	17
3919	číst	17	17
3920	čistá	17	17
3921	čisté	18	18
3922	čistě	18	18
3923	čistého	18	18
3924	čistém	18	18
3925	čistí	17	17
3926	čisticí	17	17
3927	čistič	17	17
3928	čističky	17	17
3929	čistil	17	17
3930	čistíren	17	17
3931	čistírny	17	17
3932	čistit	17	17
3933	čistky	17	17
3934	čisto	17	17
3935	čistota	17	17
3936	čistotě	18	18
3937	čistotou	17	17
3938	čistotu	17	17
3939	čistoty	17	17
3940	čistou	17	17
3941	čistší	18	18
3942	čistý	17	17
3943	čistých	17	17
3944	čistým	17	17
3945	čistými	17	17
3946	čišela	18	18
3947	čiší	18	18
3948	číši	18	18
3949	číšnice	18	18
3950	číšníci	18	18
3951	číšník	18	18
3952	číšníka	18	18
3953	čištění	18	18
3954	čištěním	18	18
3955	čítá	17	17
3956	čítající	17	17
3957	čítala	17	17
3958	čitelná	17	17
3959	čitelné	18	18
3960	čitelnost	17	17
3961	čitelný	17	17
3962	čížek	17	17
3963	čkd	17	17
3964	článcích	17	17
3965	článek	17	17
3966	článkem	17	17
3967	článku	17	17
3968	článků	17	17
3969	články	17	17
3970	člen	17	17
3971	člena	17	17
3972	členech	17	17
3973	členem	17	17
3974	členění	18	18
3975	člení	17	17
3976	členité	18	18
3977	členka	17	17
3978	členkou	17	17
3979	členky	17	17
3980	členovců	17	17
3981	členové	18	18
3982	členská	17	17
3983	členské	18	18
3984	členského	18	18
3985	členskou	17	17
3986	členských	17	17
3987	členským	17	17
3988	členskými	17	17
3989	členství	17	17
3990	členu	17	17
3991	členů	17	17
3992	členům	17	17
3993	členy	17	17
3994	člk	17	17
3995	člověče	18	18
3996	člověk	18	18
3997	člověka	18	18
3998	člověkem	18	18
3999	člověku	18	18
4000	človíček	17	17
4001	čls	17	17
4002	člun	17	17
4003	člunem	17	17
4004	člunu	17	17
4005	člunů	17	17
4006	čluny	17	17
4007	čmfs	17	17
4008	čnb	17	17
4009	ční	17	17
4010	čobe	17	17
4011	čoček	17	17
4012	čočka	17	17
4013	čočkou	17	17
4014	čočku	17	17
4015	čočky	17	17
4016	čoi	17	17
4017	čokoláda	17	17
4018	čokoládě	18	18
4019	čokoládou	17	17
4020	čokoládová	17	17
4021	čokoládové	18	18
4022	čokoládovou	17	17
4023	čokoládový	17	17
4024	čokoládových	17	17
4025	čokoládu	17	17
4026	čokolády	17	17
4027	čon	17	17
4028	čos	17	17
4029	čov	17	17
4030	čro	17	17
4031	čsa	17	17
4032	čsad	17	17
4033	čsav	17	17
4034	čsl	17	17
4035	čsn	17	17
4036	čsob	17	17
4037	čsov	17	17
4038	čsr	17	17
4039	čssd	17	17
4040	čssr	17	17
4041	čstv	17	17
4042	čsú	17	17
4043	čte	17	17
4044	čtecí	17	17
4045	čteček	17	17
4046	čtečka	17	17
4047	čtečkou	17	17
4048	čtečku	17	17
4049	čtečky	17	17
4050	čteme	17	17
4051	čtenář	17	17
4052	čtenáře	17	17
4053	čtenářem	17	17
4054	čtenáři	17	17
4055	čtenářka	17	17
4056	čtenářky	17	17
4057	čtenářské	18	18
4058	čtenářský	17	17
4059	čtenářů	17	17
4060	čtenářům	17	17
4061	čtení	17	17
4062	čtením	17	17
4063	čteš	18	18
4064	čtete	17	17
4065	čtěte	18	18
4066	čti	17	17
4067	čtk	17	17
4068	čtou	17	17
4069	čtrnáct	17	17
4070	čtrnácté	18	18
4071	čtrnáctého	18	18
4072	čtrnácti	17	17
4073	čtu	17	17
4074	čtú	17	17
4075	čtverce	17	17
4076	čtverci	17	17
4077	čtvercové	18	18
4078	čtvercový	17	17
4079	čtverců	17	17
4080	čtverec	17	17
4081	čtvereční	17	17
4082	čtverečních	17	17
4083	čtveřice	17	17
4084	čtveřici	17	17
4085	čtvrt	17	17
4086	čtvrť	19	19
4087	čtvrtá	17	17
4088	čtvrté	18	18
4089	čtvrtě	18	18
4090	čtvrteční	17	17
4091	čtvrtečním	17	17
4092	čtvrtého	18	18
4093	čtvrtek	17	17
4094	čtvrtém	18	18
4095	čtvrtfinále	17	17
4096	čtvrthodině	18	18
4097	čtvrthodinu	17	17
4098	čtvrti	17	17
4099	čtvrtí	17	17
4100	čtvrtích	17	17
4101	čtvrtin	17	17
4102	čtvrtina	17	17
4103	čtvrtině	18	18
4104	čtvrtinu	17	17
4105	čtvrtiny	17	17
4106	čtvrtka	17	17
4107	čtvrtku	17	17
4108	čtvrtky	17	17
4109	čtvrtletí	17	17
4110	čtvrtletní	17	17
4111	čtvrtou	17	17
4112	čtvrtstoletí	17	17
4113	čtvrtý	17	17
4114	čtvrtým	17	17
4115	čtyř	17	17
4116	čtyřech	17	17
4117	čtyřem	17	17
4118	čtyřhře	17	17
4119	čtyři	17	17
4120	čtyřiadvacet	17	17
4121	čtyřiadvaceti	17	17
4122	čtyřiadvacetiletý	17	17
4123	čtyřicátník	17	17
4124	čtyřicátých	17	17
4125	čtyřicet	17	17
4126	čtyřiceti	17	17
4127	čtyřicetiletý	17	17
4128	čtyřicítce	17	17
4129	čtyřicítku	17	17
4130	čtyřicítky	17	17
4131	čtyřikrát	17	17
4132	čtyřka	17	17
4133	čtyřku	17	17
4134	čtyřky	17	17
4135	čtyřleté	18	18
4136	čtyřletého	18	18
4137	čtyřletý	17	17
4138	čtyřlístek	17	17
4139	čtyřma	17	17
4140	čtyřmi	17	17
4141	čtyřválec	17	17
4142	čumák	17	17
4143	čumákem	17	17
4144	ďábel	19	19
4145	ďábelské	19	19
4146	ďábelský	19	19
4147	dabing	12	12
4148	dabingu	12	12
4149	ďábla	19	19
4150	ďáblem	19	19
4151	daimler	13	14
4152	dají	17	17
4153	dakar	7	7
4154	dal	4	4
4155	dál	16	16
4156	dala	4	4
4157	dálce	16	16
4158	dále	16	16
4159	dalece	12	12
4160	dálek	16	16
4161	daleka	8	8
4162	daleké	18	18
4163	dalekého	18	18
4164	dalekém	18	18
4165	daleko	9	9
4166	dalekohled	9	9
4167	dalekohledem	13	14
4168	dalekohledu	9	9
4169	dalekohledy	13	9
4170	dalekosáhlé	18	18
4171	dalekých	15	15
4172	dali	8	8
4173	dáli	16	16
4174	dalibor	11	11
4175	dalibora	11	11
4176	dalík	17	17
4177	dálka	16	16
4178	dálkové	18	18
4179	dálkově	18	18
4180	dálkového	18	18
4181	dálkový	16	16
4182	dálkových	16	16
4183	dálkovým	16	16
4184	dálku	16	16
4185	dálky	16	16
4186	dálném	18	18
4187	dálnic	16	16
4188	dálnice	16	16
4189	dálnici	16	16
4190	dálnicí	17	17
4191	dálnicích	17	17
4192	dálniční	17	17
4193	dálničního	17	17
4194	dálničních	17	17
4195	dálný	16	16
4196	dalo	9	9
4197	další	18	18
4198	dalšího	18	18
4199	dalších	18	18
4200	dalším	18	18
4201	dalšími	18	18
4202	dalšímu	18	18
4203	daly	13	6
4204	dáma	16	16
4205	dámám	16	16
4206	dámami	16	16
4207	damašku	18	18
4208	dáme	16	16
4209	dámě	18	18
4210	damková	16	16
4211	dámou	16	16
4212	dámské	18	18
4213	dámského	18	18
4214	dámskou	16	16
4215	dámský	16	16
4216	dámských	16	16
4217	dámu	16	16
4218	dámy	16	16
4219	dan	12	12
4220	daň	19	19
4221	dán	16	16
4222	dana	12	12
4223	daná	16	16
4224	dána	16	16
4225	daně	18	18
4226	daného	18	18
4227	daněk	18	18
4228	danem	13	14
4229	daném	18	18
4230	daněmi	18	18
4231	danému	18	18
4232	daniel	12	12
4233	daniela	12	12
4234	daniele	12	12
4235	danielem	13	14
4236	danieli	12	12
4237	danielovi	12	12
4238	danielu	12	12
4239	daních	17	17
4240	daňová	19	19
4241	daňové	19	19
4242	dánové	18	18
4243	daňového	19	19
4244	daňovém	19	19
4245	daňoví	19	19
4246	daňovou	19	19
4247	daňový	19	19
4248	daňových	19	19
4249	daňovým	19	19
4250	daňovými	19	19
4251	dánska	16	16
4252	dánská	16	16
4253	dánské	18	18
4254	dánského	18	18
4255	dánsko	16	16
4256	dánsku	16	16
4257	dánský	16	16
4258	dánských	16	16
4259	danu	12	12
4260	dany	13	12
4261	daný	15	15
4262	dány	16	16
4263	daných	15	15
4264	daným	15	15
4265	danými	15	15
4266	dar	7	7
4267	dara	7	7
4268	darby	13	11
4269	dárce	16	16
4270	dárci	16	16
4271	dárcovství	17	17
4272	dárců	16	16
4273	darebák	16	16
4274	dáreček	17	17
4275	dárečky	17	17
4276	dárek	16	16
4277	darem	13	14
4278	daria	8	8
4279	darius	8	8
4280	dark	7	7
4281	dárkem	16	16
4282	darkmarket	13	14
4283	darkmarketu	13	14
4284	dárkové	18	18
4285	dárku	16	16
4286	dárků	16	16
4287	dárky	16	16
4288	daroval	11	11
4289	darovala	11	11
4290	darovali	11	11
4291	darování	17	17
4292	darovat	11	11
4293	daru	7	7
4294	darů	7	7
4295	daruje	8	8
4296	darwin	12	12
4297	darwina	12	12
4298	dary	13	7
4299	daří	17	17
4300	dařilo	16	16
4301	dařit	16	16
4302	dásně	18	18
4303	dásní	17	17
4304	dáš	18	18
4305	dáša	18	18
4306	dat	5	6
4307	dát	16	16
4308	data	5	6
4309	databáze	16	16
4310	databázi	16	16
4311	databází	17	17
4312	databázích	17	17
4313	dáte	16	16
4314	datech	12	12
4315	datem	13	14
4316	datová	16	16
4317	datování	17	17
4318	datovat	11	11
4319	datové	18	18
4320	datového	18	18
4321	datovém	18	18
4322	datovou	11	11
4323	datový	15	15
4324	datových	15	15
4325	datovým	15	15
4326	datu	7	7
4327	datuje	8	8
4328	datují	17	17
4329	datum	13	14
4330	datům	13	14
4331	daty	13	6
4332	dav	11	11
4333	dává	16	16
4334	dávají	17	17
4335	dával	16	16
4336	dávala	16	16
4337	dávali	16	16
4338	dávalo	16	16
4339	dávaly	16	16
4340	dávám	16	16
4341	dáváme	16	16
4342	dávání	17	17
4343	dáváš	18	18
4344	dávat	16	16
4345	dáváte	16	16
4346	dávce	16	16
4347	dave	11	11
4348	davea	11	11
4349	dávej	16	16
4350	dávejte	16	16
4351	dávek	16	16
4352	davem	13	14
4353	david	11	11
4354	davida	11	11
4355	davide	11	11
4356	davidem	13	14
4357	davidova	11	11
4358	davidovi	11	11
4359	davidson	12	12
4360	davidu	11	11
4361	davis	11	11
4362	dávka	16	16
4363	dávkách	16	16
4364	dávkami	16	16
4365	dávkou	16	16
4366	dávkování	17	17
4367	dávku	16	16
4368	dávky	16	16
4369	dávná	16	16
4370	dávné	18	18
4371	dávného	18	18
4372	dávní	17	17
4373	dávno	16	16
4374	dávnou	16	16
4375	dávnověku	18	18
4376	dávný	16	16
4377	dávných	16	16
4378	dávným	16	16
4379	dávnými	16	16
4380	davosu	11	11
4381	davu	11	11
4382	davů	11	11
4383	davy	13	11
4384	dbá	16	16
4385	dbají	17	17
4386	dbal	11	11
4387	dbali	11	11
4388	dbát	16	16
4389	dbejte	11	11
4390	dcer	12	12
4391	dcera	12	12
4392	dcerám	16	16
4393	dcerami	13	14
4394	dcerce	12	12
4395	dcerka	12	12
4396	dcerkou	12	12
4397	dcerku	12	12
4398	dcerky	13	12
4399	dcerou	12	12
4400	dceru	12	12
4401	dcery	13	12
4402	dceři	16	16
4403	dceřiná	16	16
4404	dceřiné	18	18
4405	dceřinou	16	16
4406	dceřiných	16	16
4407	debakl	11	11
4408	debaklu	11	11
4409	debat	11	11
4410	debata	11	11
4411	debatách	16	16
4412	debatě	18	18
4413	debatovali	11	11
4414	debatovat	11	11
4415	debatu	11	11
4416	debaty	13	11
4417	debil	11	11
4418	debilní	17	17
4419	deborah	11	11
4420	debut	11	11
4421	debutoval	11	11
4422	debutu	11	11
4423	děcek	18	18
4424	decentně	18	18
4425	decentní	17	17
4426	deci	12	12
4427	děcka	18	18
4428	děcko	18	18
4429	děčín	18	18
4430	děčína	18	18
4431	děčíně	18	18
4432	děd	18	18
4433	děda	18	18
4434	dědeček	18	18
4435	dědečka	18	18
4436	dědečkem	18	18
4437	dědečkovi	18	18
4438	dědek	18	18
4439	dědí	18	18
4440	dědic	18	18
4441	dědice	18	18
4442	dědicem	18	18
4443	dědici	18	18
4444	dědické	18	18
4445	dědického	18	18
4446	dědicové	18	18
4447	dědictví	18	18
4448	dědictvím	18	18
4449	dědiců	18	18
4450	dědičná	18	18
4451	dědičné	18	18
4452	dědičnost	18	18
4453	dědičnosti	18	18
4454	dědičnou	18	18
4455	dědičných	18	18
4456	dědit	18	18
4457	dědka	18	18
4458	dědo	18	18
4459	dědou	18	18
4460	dědovi	18	18
4461	dědu	18	18
4462	dedukce	12	12
4463	dědy	18	18
4464	defekt	8	8
4465	defektů	8	8
4466	defekty	13	8
4467	defenzivě	18	18
4468	defenzivní	17	17
4469	deficit	12	12
4470	deficitu	12	12
4471	deficity	13	12
4472	defilé	18	18
4473	definic	12	12
4474	definice	12	12
4475	definici	12	12
4476	definicí	17	17
4477	definitivně	18	18
4478	definitivní	17	17
4479	definitivního	17	17
4480	definitivním	17	17
4481	definoval	12	12
4482	definovali	12	12
4483	definován	16	16
4484	definovaná	16	16
4485	definována	16	16
4486	definované	18	18
4487	definování	17	17
4488	definováno	16	16
4489	definovanou	12	12
4490	definovaný	15	15
4491	definovány	16	16
4492	definovaných	15	15
4493	definovat	12	12
4494	definuje	12	12
4495	definujeme	13	14
4496	definují	17	17
4497	deformace	13	14
4498	deformaci	13	14
4499	deformací	17	17
4500	degenerace	12	12
4501	degradace	12	12
4502	degradaci	12	12
4503	dech	12	12
4504	dechem	13	14
4505	dechová	16	16
4506	dechové	18	18
4507	dechovka	12	12
4508	dechovou	12	12
4509	dechových	15	15
4510	dechu	12	12
4511	dei	8	8
4512	dej	8	8
4513	děj	18	18
4514	děje	18	18
4515	dějem	18	18
4516	dějepis	18	18
4517	dějepisu	18	18
4518	ději	18	18
4519	dějí	18	18
4520	dějin	18	18
4521	dějinách	18	18
4522	dějinám	18	18
4523	dějinami	18	18
4524	dějinné	18	18
4525	dějiny	18	18
4526	dějiště	18	18
4527	dějištěm	18	18
4528	dejme	13	14
4529	dějství	18	18
4530	dejte	8	8
4531	dějů	18	18
4532	dejvicích	17	17
4533	dek	8	8
4534	deka	8	8
4535	dekád	16	16
4536	dekádách	16	16
4537	dekádě	18	18
4538	dekadentní	17	17
4539	dekádu	16	16
4540	dekády	16	16
4541	děkan	18	18
4542	děkana	18	18
4543	děkanem	18	18
4544	deklarace	12	12
4545	deklaraci	12	12
4546	deklaroval	11	11
4547	deklaruje	8	8
4548	dekolt	9	9
4549	dekor	9	9
4550	dekorace	12	12
4551	dekoraci	12	12
4552	dekorací	17	17
4553	dekorační	17	17
4554	dekorativní	17	17
4555	dekorativních	17	17
4556	dekorem	13	14
4557	dekoru	9	9
4558	dekorů	9	9
4559	dekory	13	9
4560	dekou	9	9
4561	děkoval	18	18
4562	děkovala	18	18
4563	děkovat	18	18
4564	dekret	8	8
4565	dekretu	8	8
4566	dekretů	8	8
4567	dekrety	13	8
4568	deku	8	8
4569	děkuje	18	18
4570	děkujeme	18	18
4571	děkuji	18	18
4572	děkuju	18	18
4573	deky	13	8
4574	del	8	8
4575	děl	18	18
4576	děla	18	18
4577	dělá	18	18
4578	delacorte	12	12
4579	dělaj	18	18
4580	dělají	18	18
4581	dělal	18	18
4582	dělala	18	18
4583	dělali	18	18
4584	dělalo	18	18
4585	dělaly	18	18
4586	dělám	18	18
4587	děláme	18	18
4588	děláš	18	18
4589	dělat	18	18
4590	děláte	18	18
4591	dělávají	18	18
4592	dělával	18	18
4593	dělávala	18	18
4594	dělba	18	18
4595	dělby	18	18
4596	délce	18	18
4597	déle	18	18
4598	delegace	12	12
4599	delegaci	12	12
4600	delegáti	16	16
4601	delegátů	16	16
4602	dělej	18	18
4603	dělejte	18	18
4604	délek	18	18
4605	dělem	18	18
4606	dělení	18	18
4607	dělením	18	18
4608	delfíni	17	17
4609	delfínů	17	17
4610	dělí	18	18
4611	dělicí	18	18
4612	dělící	18	18
4613	delikátní	17	17
4614	delikt	8	8
4615	deliktů	8	8
4616	delikty	13	8
4617	dělil	18	18
4618	dělila	18	18
4619	dělili	18	18
4620	dělilo	18	18
4621	dělily	18	18
4622	dělíme	18	18
4623	dělit	18	18
4624	délka	18	18
4625	délkách	18	18
4626	délkou	18	18
4627	délku	18	18
4628	délky	18	18
4629	dělnic	18	18
4630	dělnice	18	18
4631	dělníci	18	18
4632	dělnické	18	18
4633	dělnických	18	18
4634	dělník	18	18
4635	dělníka	18	18
4636	dělníků	18	18
4637	dělníkům	18	18
4638	dělníky	18	18
4639	dělo	18	18
4640	dělohu	18	18
4641	dělohy	18	18
4642	deloitte	9	9
4643	dělostřelecké	18	18
4644	dělostřelectva	18	18
4645	dělostřelectvo	18	18
4646	děloze	18	18
4647	děložní	18	18
4648	děložního	18	18
4649	delší	18	18
4650	delšího	18	18
4651	delších	18	18
4652	delším	18	18
4653	delšími	18	18
4654	delta	8	8
4655	deltě	18	18
4656	delty	13	8
4657	děly	18	18
4658	dem	13	14
4659	demence	13	14
4660	demencí	17	17
4661	demirčjan	17	17
4662	demisi	13	14
4663	demografické	18	18
4664	demografický	15	15
4665	demokracie	13	14
4666	demokracii	13	14
4667	demokracií	17	17
4668	demokrat	13	14
4669	demokraté	18	18
4670	demokratická	16	16
4671	demokratické	18	18
4672	demokratického	18	18
4673	demokratickém	18	18
4674	demokratickou	13	14
4675	demokraticky	13	14
4676	demokratický	15	15
4677	demokratických	15	15
4678	demokratů	13	14
4679	demokratům	13	14
4680	demokraty	13	14
4681	demolice	13	14
4682	demolici	13	14
4683	démon	18	18
4684	démona	18	18
4685	démoni	18	18
4686	demonstrace	13	14
4687	demonstraci	13	14
4688	demonstrací	17	17
4689	demonstracích	17	17
4690	demonstranti	13	14
4691	demonstrantů	13	14
4692	demonstranty	13	14
4693	demonstrativně	18	18
4694	demonstroval	13	14
4695	demonstrovat	13	14
4696	demonstruje	13	14
4697	demontáž	16	16
4698	démonů	18	18
4699	démony	18	18
4700	den	12	12
4701	dění	18	18
4702	denících	17	17
4703	denik	12	12
4704	deník	17	17
4705	deníkem	17	17
4706	deníku	17	17
4707	deníků	17	17
4708	deníky	17	17
4709	děním	18	18
4710	denis	12	12
4711	denisa	12	12
4712	denise	12	12
4713	denně	18	18
4714	denní	17	17
4715	denního	17	17
4716	denních	17	17
4717	denním	17	17
4718	dennis	12	12
4719	dennodenně	18	18
4720	dentální	17	17
4721	denveru	12	12
4722	depeše	18	18
4723	depozitáře	16	16
4724	deprese	10	10
4725	depresemi	13	14
4726	depresi	10	10
4727	depresí	17	17
4728	depresivní	17	17
4729	deprimující	17	17
4730	děr	18	18
4731	derby	13	11
4732	derivace	12	12
4733	deriváty	16	16
4734	derou	9	9
4735	děs	18	18
4736	desátá	16	16
4737	desáté	18	18
4738	desátého	18	18
4739	desátém	18	18
4740	desátou	16	16
4741	desátý	16	16
4742	descartes	12	12
4743	desce	12	12
4744	desek	8	8
4745	děsem	18	18
4746	deset	8	8
4747	deseti	8	8
4748	desetiletá	16	16
4749	desetileté	18	18
4750	desetiletého	18	18
4751	desetiletí	17	17
4752	desetiletích	17	17
4753	desetiletou	9	9
4754	desetiletý	15	15
4755	desetin	12	12
4756	desetina	12	12
4757	desetinu	12	12
4758	desetiny	13	12
4759	desetiprocentní	17	17
4760	desetitisíce	17	17
4761	desetkrát	16	16
4762	děsí	18	18
4763	design	12	12
4764	designbloku	12	12
4765	designem	13	14
4766	designér	18	18
4767	designéra	18	18
4768	designérka	18	18
4769	designérů	18	18
4770	designéry	18	18
4771	designéři	18	18
4772	designové	18	18
4773	designově	18	18
4774	designový	15	15
4775	designových	15	15
4776	designu	12	12
4777	děsil	18	18
4778	děsila	18	18
4779	děsilo	18	18
4780	desire	8	8
4781	děsit	18	18
4782	desítce	17	17
4783	desítek	17	17
4784	desíti	17	17
4785	desítka	17	17
4786	desítkách	17	17
4787	desítkám	17	17
4788	desítkami	17	17
4789	desítku	17	17
4790	desítky	17	17
4791	děsivá	18	18
4792	děsivé	18	18
4793	děsivě	18	18
4794	děsivého	18	18
4795	děsivou	18	18
4796	děsivý	18	18
4797	děsivých	18	18
4798	deska	8	8
4799	deskách	16	16
4800	deskami	13	14
4801	deskou	9	9
4802	deskové	18	18
4803	deskriptivní	17	17
4804	desku	8	8
4805	desky	13	8
4806	děsná	18	18
4807	děsně	18	18
4808	děsný	18	18
4809	despektem	13	14
4810	destiček	17	17
4811	destičku	17	17
4812	destičky	17	17
4813	destinace	12	12
4814	destinaci	12	12
4815	destinací	17	17
4816	destrukce	12	12
4817	destrukci	12	12
4818	destruktivní	17	17
4819	děsu	18	18
4820	dešifrovat	18	18
4821	déšť	19	19
4822	deště	18	18
4823	deštěm	18	18
4824	dešti	18	18
4825	deštích	18	18
4826	deštivé	18	18
4827	deštivého	18	18
4828	deštného	18	18
4829	deštník	18	18
4830	deštníkem	18	18
4831	deštníku	18	18
4832	deštníky	18	18
4833	deštných	18	18
4834	dešťová	19	19
4835	dešťové	19	19
4836	dešťovou	19	19
4837	dešťových	19	19
4838	dešťů	19	19
4839	detail	8	8
4840	detailech	12	12
4841	detailem	13	14
4842	detailně	18	18
4843	detailněji	18	18
4844	detailnější	18	18
4845	detailní	17	17
4846	detailních	17	17
4847	detailu	8	8
4848	detailů	8	8
4849	detaily	13	8
4850	děťátka	19	19
4851	děťátko	19	19
4852	dětech	18	18
4853	detekce	12	12
4854	detekci	12	12
4855	detekovat	11	11
4856	detektiv	11	11
4857	detektiva	11	11
4858	detektive	11	11
4859	detektivek	11	11
4860	detektivka	11	11
4861	detektivku	11	11
4862	detektivky	13	11
4863	detektivní	17	17
4864	detektivové	18	18
4865	detektivů	11	11
4866	detektivy	13	11
4867	detektor	9	9
4868	detektorem	13	14
4869	detektoru	9	9
4870	detektorů	9	9
4871	detektory	13	9
4872	dětem	18	18
4873	děti	18	18
4874	dětí	18	18
4875	dětinské	18	18
4876	dětma	18	18
4877	dětmi	18	18
4878	detroit	9	9
4879	detroitu	9	9
4880	dětská	18	18
4881	dětské	18	18
4882	dětského	18	18
4883	dětském	18	18
4884	dětskému	18	18
4885	dětskou	18	18
4886	dětsky	18	18
4887	dětský	18	18
4888	dětských	18	18
4889	dětským	18	18
4890	dětskými	18	18
4891	dětství	18	18
4892	deutsche	12	12
4893	deutschland	12	12
4894	devadesát	16	16
4895	devadesáti	16	16
4896	devadesátých	16	16
4897	devastaci	12	12
4898	devátá	16	16
4899	deváté	18	18
4900	devátého	18	18
4901	devátém	18	18
4902	devatenáct	16	16
4903	devatenáctého	18	18
4904	devatenáctém	18	18
4905	devatenácti	16	16
4906	devatenáctiletý	16	16
4907	devátou	16	16
4908	devátý	16	16
4909	děvčat	18	18
4910	děvčata	18	18
4911	děvčátka	18	18
4912	děvčátko	18	18
4913	děvčatům	18	18
4914	děvčaty	18	18
4915	děvče	18	18
4916	děvčete	18	18
4917	developer	11	11
4918	developera	11	11
4919	developerské	18	18
4920	developerů	11	11
4921	developeři	16	16
4922	development	13	14
4923	děvenko	18	18
4924	devět	18	18
4925	devětkrát	18	18
4926	devils	11	11
4927	devíti	17	17
4928	děvka	18	18
4929	děvky	18	18
4930	dex	13	14
4931	dezert	8	14
4932	dezerty	13	14
4933	dezinfekce	12	14
4934	dezinfekci	12	14
4935	dezinfekční	17	17
4936	dhark	7	7
4937	diabetem	13	14
4938	diabetes	11	11
4939	diabetiků	11	11
4940	diabetu	11	11
4941	diagnostice	12	12
4942	diagnostické	18	18
4943	diagnostických	15	15
4944	diagnostika	12	12
4945	diagnostikovali	12	12
4946	diagnostikovat	12	12
4947	diagnostiku	12	12
4948	diagnostiky	13	12
4949	diagnóza	19	19
4950	diagnóze	19	19
4951	diagnózou	19	19
4952	diagnózu	19	19
4953	diagnózy	19	19
4954	diagram	13	14
4955	diagramu	13	14
4956	diagramů	13	14
4957	diagramy	13	14
4958	dialog	9	9
4959	dialogu	9	9
4960	dialogů	9	9
4961	dialogy	13	9
4962	diamant	13	14
4963	diamantem	13	14
4964	diamantové	18	18
4965	diamantu	13	14
4966	diamantů	13	14
4967	diamanty	13	14
4968	diametrálně	18	18
4969	diamond	13	14
4970	diana	12	12
4971	diane	12	12
4972	dianě	18	18
4973	dianu	12	12
4974	diany	13	12
4975	diář	16	16
4976	diáře	16	16
4977	dick	12	12
4978	didaktické	18	18
4979	didaktických	15	15
4980	didaktika	8	8
4981	didaktiky	13	8
4982	die	8	8
4983	diecéze	18	18
4984	diega	8	8
4985	diego	9	9
4986	diegu	8	8
4987	dienstbier	12	12
4988	dienstbiera	12	12
4989	diesel	8	8
4990	diet	8	8
4991	dieta	8	8
4992	dietě	18	18
4993	dieter	8	8
4994	dietmar	13	14
4995	dietní	17	17
4996	dietou	9	9
4997	dietrich	12	12
4998	dietu	8	8
4999	diety	13	8
5000	diference	12	12
\.
-- 5000 snad bude stacit

COPY public.texty (id, jmeno, text1, text2, text3, text4, text5, text6, text7, text8, text9, text10) FROM stdin;
1	Farma zvířat	Pan Jones zamkl sice na Panské farmě na noc kurníky, ale byl tak\nzpitý, že zapomněl zastrčit závlačky na dveřích. Ve světle sem tam se\npohupující lucerny dovrávoral přes dvůr k zadním dveřím, skopl\nholiny, natočil si poslední sklenici piva ze sudu v kuchyňce a dopotácel se do postele, ve které již nějakou chvíli pochrupovala paní\nJonesová.	Člověk je jediný tvor, který konzumuje, aniž by produkoval. Nedává mléko, nesnáší vejce, je příliš slabý na to, aby tahal pluh, nedokáže utíkat tak rychle, aby chytil králíka! A přesto je pánem všech\nzvířat. Nutí je pracovat, za odvedenou dřinu se jim odvděčuje jen\ntolik, aby nepošla hlady a mohla pro něj dál lopotit, a zbytek si nechává pro sebe.	Tam, u dlouhého stolu, sedělo půl tuctu farmářů a půl tuctu nejprominentnějších prasat. Napoleon sám zaujímal čestné místo v čele stolu. Ostatní prasata se rozvalovala ve svých křeslech. Společnost se očividně dobře bavila při partii karet, která byla však nyní přerušena, jelikož se chystal přípitek.	Zima byla krutá. Bouřky střídaly plískanice a sníh, pak udeřily třeskuté mrazy, které polevily až koncem února. Ale zvířata se na přestavbě Větrného mlýna činila, stůj co stůj, neboť dobře věděla, že je\nokolní svět bedlivě sleduje a že by závistiví lidé skákali radostí, kdyby\nse farmě nepodařilo dílo včas dokončit.	Jednoho nedělního rána, když se zvířata shromáždila, aby vyslechla rozkazy, oznámil Napoleon, že se rozhodl nastolit novou politiku. Odteď začne farma navazovat obchodní styky se sousedními\nfarmami: samozřejmě nepůjde o žádné obchodování, cílem bude\npouhé získání nezbytně nutného materiálu. 	Rámus totiž probudil pana Jonese, který byl skokem venku z postele, přesvědčen, že se na dvůr vloudila liška. Popadl pušku, kterou měl vždy připravenou, a šestkrát vypálil z okna do tmy. Kulky se zaryly do zdi stodoly a shromáždění se chvatně rozuteklo. Ptáci vzlétli na bidýlka, zvířata ulehla na slámu a chvíli nato se celá farma ponořila do spánku.	Druhý den se probudili za úsvitu jako obvykle. Vtom si vzpomněli na tu slávu sláv, která se udála předchozího dne, a vyrazili jako o závod společně na pastvu. Kousek od ní byl pahorek skýtající výhled na celou farmu. Vyběhli na něj a z vrcholu se pak v třpytivě jasném ranním světle nemohli vynadívat na vše, co viděli kolem.	Vrátili se všichni k hospodářským stavením, kam nechali Sněžník s Napoleonem donést žebřík a opřít ho o vrata stodoly. Vysvětlili ostatním zvířatům, že se jim za poslední tři měsíce studia podařilo zredukovat principy animalismu na Sedmero přikázání. Těchto sedm nezměnitelných zákonů bude nyní napsáno na vrata a všechna zvířata na farmě podle nich musí již navždy žít.	Prasata však byla natolik chytrá, že dokázala najít cestičku z každé zapeklité situace. Koně zase znali pole jako svá kopyta, a rozuměli tak všemu, co se točilo kolem sekání a hrabání sena, lépe než pan Jones a jeho lidi. Prasata, po pravdě řečeno, sice nepracovala, zato veškeré práce rozdělovala, řídila a dohlížela na jejich vykonávání.	I když krmení nebylo víc než za časů pana Jonese, nebylo ho ani míň. Fakt, že živila akorát sama sebe, a ne pět lidských darmožroutů navíc, skýtal takovou výhodu, že by zvířata bývala musela dělat opravdu velké chyby, aby byla ve ztrátě. Dalším jejich velkým plusem bylo, že jejich pracovní metody byly efektivnější a produktivnější.
2	Malý princ	Skutečně. Když je v Americe poledne, víme, že nad Francií slunce zapadá. Stačilo by, abychom se mohli ocitnout za minutu ve Francii, a viděli bychom západ slunce. Bohužel Francie je příliš daleko. Ale na tvé malinké planetě ti stačilo posunout židli o několik kroků. A díval ses na soumrak, kdykoliv se ti zachtělo...	Je to těžké, v mém věku se dát znovu do kreslení, když se člověk nepokusil o nic jiného než v šesti letech o zavřeného a otevřeného hroznýše! Pokusím se ovšem nakreslit portréty co nejvěrněji. Ale nejsem si tak docela jistý, zda se mi to podaří. Jedna kresba se zdaří, a druhá už ne.	To mě nemohlo příliš překvapit. Dobře jsem věděl, že kromě velkých planet, jako jsou Země, Jupiter, Mars, Venuše, které dostaly jméno, existují ještě stovky jiných, které jsou někdy tak malé, že dá hodně práce spatřit je aspoň dalekohledem. Když hvězdář takovou planetu objeví, dá jí místo jména číslo.	Má kresba ovšem není zdaleka tak půvabná jako model. Ale za to já nemohu. Dospělí mě odradili od malířské kariéry, když mi bylo šest let, a proto jsem se nenaučil kreslit nic jiného než zavřené a otevřené hroznýše. Udiveně jsem se díval na ten zjev. Považte jen, že jsem byl na tisíc mil od jakéhokoliv obydleného kraje.	Tak jsem žil sám a neměl jsem nikoho, s kým bych si mohl opravdu popovídat. Tu se mi jednou před šesti lety v poušti na Sahaře porouchal motor. Něco se v něm polámalo. A poněvadž jsem neměl s sebou mechanika ani cestující, chtěl jsem se do té nesnadné opravy pustit sám. Byla to pro mne otázka života nebo smrti.	Možná je ten člověk zbytečný. A přece je méně zbytečný než král, než domýšlivec, než byznysmen a než pijan. Jeho práce má alespoň smysl. Když rozsvítí svítilnu, jako by se zrodilo o hvězdu nebo o květinu víc. Když zhasne květinu, jako by květina nebo hvězda šly spát. Je to moc hezké zaměstnání. A opravdu užitečné, protože je hezké.	Na toho, kdo viděl Zemi trochu z dálky, dělalo to ohromný dojem. Pohyby této armády byly řízeny jako operní balet. Nejprve přišli na řadu lampáři na Novém Zélandě a v Austrálii. A ti, když rozžali lampy, šli spát. Po nich nastoupili do toho reje lampáři v Číně a na Sibiři. Potom také oni obratně zmizeli za kulisami. Tu nastoupili lampáři ruští a indičtí.	Nebyl jsem moc poctivý, když jsem vám vyprávěl o tom lampáři. Riskuji, že vyvolám špatnou představu o naší planetě u těch, kdo ji neznají. Lidé zabírají na Zemi velice málo místa. Kdyby ty dvě miliardy obyvatel, kteří zalidňují Zemi, stály trochu stlačeny jako na táboru lidu, vešly by se snadno na náměstí dvacet mil dlouhé a dvacet mil široké.	Všechny slepice jsou si navzájem podobné a také lidé jsou si podobní. Trochu se proto nudím. Ale když si mě ochočíš, bude můj život jakoby prozářen sluncem. Poznám zvuk kroků, který bude jiný než všechny ostatní. Ostatní kroky mě zahánějí pod zem. Ale tvůj krok mě jako hudba vyláká z doupěte. A pak, podívej se! Vidíš tamhleta obilná pole?	Byl jsem překvapen, že pojednou chápu to tajemné záření písku. Když jsem byl malým chlapcem, bydlil jsem ve starobylém domě a pověst vyprávěla, že je tam zakopán poklad. Nikdy jej ovšem nikdo nedovedl objevit a snad jej ani nehledal. Ale dodával kouzlo celému tomu domu. Můj dům skrýval ve svých hlubinách tajemství...
3	O pejskovi a kočičce	Pejsek vlezl do vany a kočička ho vyprala. Drhla ho na té valše tak silně, že ji pejsek prosil, ať tolik netlačí, že by se mu mohly do sebe zamotat nohy. Když byl pejsek umytý, vlezla zas do necek kočička a pejsek ji vypral a tlačil tak silně, že ho prosila, aby ji na té valše tolik nedřel, že jí vydře do kožichu díru.	Tak do toho svého dortu dávali a míchali všechno možné, dali tam i česnek a pepř a namíchali tam sádlo i bonbóny, škvarky a skořici, krupičnou kaši a tvaroh, perník a ocet, kakao a zelí, jednu hlavu z husy a hrozinky, inu všechno možné do toho dortu dali, jen chleba tam nedali, protože pejskové a kočičky chleba zrovna tuze moc rádi nejedí.	A když spala, zavolal ten andělíček ještě jiné andělíčky a ti tu košilku krásně vyšili a udělali na ní proužečky a červené tečky a čtverečky a ještě takové všelijaké věci, jak to má kdo rád, a ještě ji pěkně vyžehlili a pak zas odletěli. A když se ráno ta noční košilka probudila a oblékla se, tu se ani nemohla poznat.	Já vám nejdřív plakala, pak jsem vrčela, nato skuhrala a nakonec syčela, a když už jsem si nevěděla rady, tak jsem ji ťapičkou sekla. Jenom trochu, ale to svět neviděl, co ta holka nadělala křiku! Tak byste měl, pane Čapek, do těch vánočních novin napsat, aby děti netahaly nás kočky za ocas, nebo si nikdy s nimi nebudeme hrát.	Když děti chtějí udělat haf!, to ti přitom nadělají nějakých chyb! Správně hafnout, to není, holenku, jen tak. To ti to musí nejdřív v prsou zarachotit a pak to musí z tebe vyštěknout rychle a statečně, jako když vystřelí, musíš při tom krásně hodit hlavou a škubnout hřbetem a zadníma nohama.	Pejsek si narovnal uši a šli. Cestou si vykládali, co budou v lese dělat. Že si tam budou hrát na schovávanou a na Jeníčka a na Mařenku a chvíli taky na honěnou. A pak že si lehnou na záda do trávy a budou se dívat nahoru na modré nebe. Co tak šli, koukal se na ně z křoví zajíc.	Kočička svlékla panenku z mokrých šatečků, vyprala je, a dala je na sluníčko usušit, aby zas byly čisté a suché. Pejsek uložil panenku do postýlky, aby se zahřála, a přinesl jí rohlíček a hrneček mléka. Panenka se najedla a pak v postýlce pěkně usnula. Pejsek s kočičkou chodili po špičkách kolem, aby panenku neprobudili, a těšili se z toho, že teď mají svoje děťátko.	A to všechno postavili kolem postýlky té pohozené panence. Panenka se probudila a měla velikou radost, když kolem sebe viděla tolik hraček, a moc se jí to u pejska a kočičky líbilo. Hrála si a hrála, hrála si s kočičkou a pejskem se všemi těmi hračkami a bylo jich tolik, že si ani dost se všemi vyhrát nemohla.	To při mňoukání se zas musejí krásně vykulit oči, hlavička dát drobet na stranu, trošičku při tom povstat a jemně se protáhnout. To by se děti tomu tuze dlouho musely ve škole učit, než by se to správně naučily, a vidíš, my kočky to umíme hned samy od sebe, už od narození. Ať to s nimi ve škole paní učitelka zkusí a budou vidět, jak je to těžké.	Sněhu habaděj pro lyžaře i pro sáňkaře, i těch klouzaček bylo letos dost. Tak se sypal, jako když se pytel s moukou roztrhne. Pejskovi a kočičce se to docela líbilo a spokojeně ťapali s tím psaním na poštu. Donesli psaní na poštu, kočička dala pejskovi olíznout známku, přilepili ji a připlácli pacičkami a dali do schránky.
\.

-- protože nefungujou emaily tak tady je testovací uživatel heslo: testtest
INSERT INTO uzivatel (email, jmeno, heslo) VALUES ('test@test.test', 'test', '$2a$10$QOq5R8EPz3Ei/LqI4Bzpm.yPSvFzSdqxK8LK34HdL/SztNngSFXFC');