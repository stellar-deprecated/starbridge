-- +migrate Up
CREATE TABLE signature_requests (
    withdraw_chain character varying(40) NOT NULL,
    deposit_chain character varying(40) NOT NULL,
    requested_action character varying(40) NOT NULL,
    deposit_id text NOT NULL,
    PRIMARY KEY (deposit_id, deposit_chain, withdraw_chain, requested_action)
);

CREATE TABLE history_stellar_transactions (
    hash character varying(64) NOT NULL PRIMARY KEY,
    envelope text NOT NULL,
    memo_hash character varying(64) NOT NULL
);
CREATE INDEX history_stellar_transactions_memo_hash ON history_stellar_transactions USING BTREE(memo_hash);

CREATE TABLE outgoing_stellar_transactions (
    envelope text NOT NULL,
    sequence bigint NOT NULL ,
    source_account text NOT NULL,
    requested_action character varying(40) NOT NULL,
    deposit_id text NOT NULL
);
CREATE UNIQUE INDEX outgoing_stellar_transaction_for_action ON outgoing_stellar_transactions USING BTREE(requested_action, deposit_id);

CREATE TABLE ethereum_signatures (
   address TEXT NOT NULL,
   token TEXT NOT NULL,
   amount TEXT NOT NULL,
   signature TEXT NOT NULL,
   expiration BIGINT NOT NULL ,
   requested_action character varying(40) NOT NULL,
   deposit_id TEXT NOT NULL
);
CREATE UNIQUE INDEX ethereum_signatures_for_action ON ethereum_signatures USING BTREE(requested_action, deposit_id);

CREATE TABLE okx_signatures (
   address TEXT NOT NULL,
   token TEXT NOT NULL,
   amount TEXT NOT NULL,
   signature TEXT NOT NULL,
   expiration BIGINT NOT NULL ,
   requested_action character varying(40) NOT NULL,
   deposit_id TEXT NOT NULL
);
CREATE UNIQUE INDEX okx_signatures_for_action ON okx_signatures USING BTREE(requested_action, deposit_id);

CREATE TABLE concordium_signatures (
   address TEXT NOT NULL,
   token TEXT NOT NULL,
   amount TEXT NOT NULL,
   signature TEXT NOT NULL,
   expiration BIGINT NOT NULL ,
   requested_action character varying(40) NOT NULL,
   deposit_id TEXT NOT NULL
);
CREATE UNIQUE INDEX concordium_signatures_for_action ON concordium_signatures USING BTREE(requested_action, deposit_id);

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

CREATE TABLE okx_deposits (
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

CREATE TABLE concordium_deposits (
    id TEXT NOT NULL PRIMARY KEY,
    amount TEXT NOT NULL,
    destination TEXT NOT NULL,
    sender TEXT NOT NULL,
    block_hash TEXT NOT NULL,
    block_time BIGINT NOT NULL
);

CREATE TABLE stellar_deposits (
   id TEXT NOT NULL PRIMARY KEY,
   ledger_time BIGINT NOT NULL,
   amount TEXT NOT NULL,
   destination TEXT NOT NULL,
   sender TEXT NOT NULL,
   asset TEXT NOT NULL
);

CREATE TABLE key_value_store (
  key varchar(255) NOT NULL,
  value varchar(255) NOT NULL,
  PRIMARY KEY (key)
);

-- +migrate Down
drop table key_value_store cascade;
drop table stellar_deposits cascade;
drop table ethereum_deposits cascade;
drop table ethereum_signatures cascade;
drop table okx_deposits cascade;
drop table okx_signatures cascade;
drop table concordium_deposits cascade;
drop table outgoing_stellar_transactions cascade;
drop table history_stellar_transactions cascade;
drop table signature_requests cascade;