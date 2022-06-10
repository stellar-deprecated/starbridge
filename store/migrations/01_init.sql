-- +migrate Up
CREATE TABLE signature_requests (
    incoming_type character varying(40) NOT NULL,
    incoming_transaction_hash text NOT NULL
);

CREATE UNIQUE INDEX type_hash ON signature_requests USING BTREE(incoming_type, incoming_transaction_hash);

CREATE TABLE outgoing_stellar_transactions (
    state character varying(20) NOT NULL,
    hash character varying(64) NOT NULL PRIMARY KEY,
    envelope text NOT NULL,
    expiration timestamp without time zone NOT NULL,
    incoming_type character varying(40) NOT NULL,
    incoming_transaction_hash text NOT NULL
);

CREATE UNIQUE INDEX outgoing_stellar_type_hash ON outgoing_stellar_transactions USING BTREE(incoming_type, incoming_transaction_hash);

CREATE TABLE incoming_ethereum_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    value_wei bigint NOT NULL,
    stellar_address character varying(56) NOT NULL
);

-- +migrate Down
drop table incoming_ethereum_transactions cascade;
drop table outgoing_stellar_transactions cascade;
drop table signature_requests cascade;