-- Table: notfy.email

-- DROP TABLE notfy.email;

CREATE TABLE notfy.email
(
    email_id integer NOT NULL,
    status_id integer NOT NULL,
    "from" character varying(10000) COLLATE pg_catalog."default" NOT NULL,
    "to" character varying(10000)[] COLLATE pg_catalog."default",
    cc character varying(10000)[] COLLATE pg_catalog."default",
    bcc character varying(10000)[] COLLATE pg_catalog."default",
    subject character varying(1000) COLLATE pg_catalog."default",
    body text COLLATE pg_catalog."default",
    CONSTRAINT email_pkey PRIMARY KEY (email_id),
    CONSTRAINT email_status FOREIGN KEY (email_id)
        REFERENCES notfy.status (status_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE notfy.email
    OWNER to postgres;
