--
-- PostgreSQL database dump
--

\restrict 4eO0PcTkjiC05jseYGoxlIUbYH4mcLfCGavNTsLm8DSCMWmBJkwtJcjUYRdu8L8

-- Dumped from database version 18.1 (Ubuntu 18.1-1.pgdg24.04+2)
-- Dumped by pg_dump version 18.1 (Ubuntu 18.1-1.pgdg24.04+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: houses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.houses (
    id integer NOT NULL,
    name character varying(200) NOT NULL
);


ALTER TABLE public.houses OWNER TO postgres;

--
-- Name: houses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.houses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.houses_id_seq OWNER TO postgres;

--
-- Name: houses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.houses_id_seq OWNED BY public.houses.id;


--
-- Name: houses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.houses ALTER COLUMN id SET DEFAULT nextval('public.houses_id_seq'::regclass);


--
-- Data for Name: houses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.houses (id, name) FROM stdin;
3	Targaryen
4	Greyjoy
5	Tully
6	Lannister
7	Arryn
8	Tyrell
9	Highgarden
10	XH6ZGZ7H2OOM52F5YBZPDGPY25
11	NNQ74WTFEV5K6HZV4JP5EE5BXS
12	FSKJ4ZDKG34CSNQD6N54NWW45G
13	GCLGEDDULZJC3LAZI4NH7M2ZYO
14	XXQWQTHRYN26UBNRPST4VYGLJW
15	KBEJEUEHCBZGYSGRLDT47NX3DY
16	MIMCKA2SN52QO24SVBV3HZNL36
17	JDY2LNHOPO567Y3P4EOC7BZOU6
\.


--
-- Name: houses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.houses_id_seq', 17, true);


--
-- Name: houses houses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.houses
    ADD CONSTRAINT houses_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict 4eO0PcTkjiC05jseYGoxlIUbYH4mcLfCGavNTsLm8DSCMWmBJkwtJcjUYRdu8L8

