CREATE TABLE notfy.status
(
    status_id bigserial NOT NULL,
    code smallint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    last_update_at timestamp without time zone NOT NULL,
    PRIMARY KEY (status_id)
)
WITH (
    OIDS = FALSE
);

ALTER TABLE notfy.status
    OWNER to postgres;
