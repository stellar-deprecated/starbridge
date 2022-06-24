-- +migrate Up
CREATE TABLE signature_requests (
    incoming_type character varying(40) NOT NULL,
    incoming_transaction_hash text NOT NULL
);

CREATE UNIQUE INDEX type_hash ON signature_requests USING BTREE(incoming_type, incoming_transaction_hash);

CREATE TABLE history_stellar_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    envelope text NOT NULL,
    memo_hash character varying(64) NOT NULL
);

CREATE INDEX history_stellar_transactions_memo_hash ON history_stellar_transactions USING BTREE(memo_hash);

CREATE TABLE outgoing_stellar_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    state character varying(20) NOT NULL,
    envelope text NOT NULL,
    incoming_type character varying(40) NOT NULL,
    incoming_transaction_hash text NOT NULL
);

CREATE UNIQUE INDEX outgoing_stellar_type_hash ON outgoing_stellar_transactions USING BTREE(incoming_type, incoming_transaction_hash);

CREATE TABLE incoming_ethereum_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    value_wei text NOT NULL,
    stellar_address character varying(56) NOT NULL,
    withdraw_expiration timestamp without time zone NOT NULL,
    withdrawn boolean
);

CREATE TABLE key_value_store (
  key varchar(255) NOT NULL,
  value varchar(255) NOT NULL,
  PRIMARY KEY (key)
);

-- +migrate Down
drop table key_value_store cascade;
drop table incoming_ethereum_transactions cascade;
drop table outgoing_stellar_transactions cascade;
drop table signature_requests cascade;