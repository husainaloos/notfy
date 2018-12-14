CREATE TABLE notfy.email
(
	    email_id bigserial NOT NULL,
	    "from" character varying(100),
	    "to" character varying(100)[],
	    cc character varying(100)[],
	    bcc character varying(100)[],
	    subject character varying(1000),
	    body text,
	    status_events jsonb,
	    PRIMARY KEY (email_id)
)
WITH (
	    OIDS = FALSE
);

ALTER TABLE notfy.email
    OWNER to postgres;
