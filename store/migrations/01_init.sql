-- +migrate Up
CREATE TABLE signature_requests (
    deposit_chain character varying(40) NOT NULL,
    requested_action character varying(40) NOT NULL,
    deposit_id text NOT NULL,
    PRIMARY KEY (deposit_id, deposit_chain, requested_action)
);

CREATE TABLE history_stellar_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    envelope text NOT NULL,
    memo_hash character varying(64) NOT NULL
);

CREATE INDEX history_stellar_transactions_memo_hash ON history_stellar_transactions USING BTREE(memo_hash);

CREATE TABLE outgoing_stellar_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    envelope text NOT NULL,
    sequence bigint NOT NULL ,
    requested_action character varying(40) NOT NULL,
    deposit_id text NOT NULL
);
CREATE UNIQUE INDEX outgoing_stellar_transaction_for_action ON outgoing_stellar_transactions USING BTREE(requested_action, deposit_id);

CREATE TABLE ethereum_deposits (
    id TEXT NOT NULL PRIMARY KEY,
    hash TEXT NOT NULL,
    log_index INTEGER NOT NULL,
    block_number BIGINT NOT NULL,
    block_time BIGINT NOT NULL,
    amount TEXT NOT NULL,
    destination TEXT NOT NULL,
    sender TEXT NOT NULL,
    token TEXT NOT NULL
);

CREATE TABLE key_value_store (
  key varchar(255) NOT NULL,
  value varchar(255) NOT NULL,
  PRIMARY KEY (key)
);

-- +migrate Down
drop table key_value_store cascade;
drop table ethereum_deposits cascade;
drop table outgoing_stellar_transactions cascade;
drop table history_stellar_transactions cascade;
drop table signature_requests cascade;