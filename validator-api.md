This document describes the backend API, behind which sits the core validator logic. The backend API
should abstract away details of the Horizon and Ethereum APIs as well as the user HTTP API.

The backend API clearly separates the core validator logic from more mundane data transformations like
parsing, filtering relevant information, etc. It should also simplify model-based testing of the core
validator logic using Ivy.

Components
==========

-   The Stellar observer

Consists of a Horizon instance and a validator shim that queries, filters, and processes data
delivered by Horizon.

-   The Ethereum observer

Consists of an Ethereum node or a connection to some API like Infura and a validator shim that further
filters and processes data.

-   The backend

Consists of an SQL database and a set of functions that process user requests and callbacks from the
Stellar and Ethereum observers.

The backend is purely reactive, i.e. it never initiates anything. If polling is needed, e.g. polling
Horizon, this should be driven by the Horizon shim in the Stellar observer component.

-   The user-facing shim

Processes HTTP requests made by the user and calls the backend.

APIs
====

Backend
-------

-   Provides `request_stellar_withdraw_tx(hash h, signature s)`.
    -   Preconditions:
        -   `h` is hash identifying an Ethereum transaction transferring an amount `a` of tokens to
            the bridge contract. This calls appears in a block that is final.
        -   `s` is `h` signed by the private key of the source account of the transaction identified
            by `h`.
    -   Returns:
        -   A signed Stellar `TransactionEnvelope` transferring `a` amount of tokens 

User-facing API
===============

Signature request
-----------------

Parameters: \* A hash identifying a deposit to the bridge contract on Ethereum

Refund request
--------------

Stellar-facing API
==================

Callbacks
---------

-   New ledger notification
    -   Must include:
        -   ledger close time
        -   all deposits made to the bridge account in this ledger
        -   all withdrawals from the bridge account made in this ledger
    -   Precondition: ledger is final

Calls
-----

-   Get the sequence number of an account
-   Ingest history (used on startup) Provides the history of all withdrawals from the bridge account

Ethereum-facing API
===================

-   New-block notification
    -   Must include:
        -   block timestamp
        -   all deposits made to the bridge contract in this block
        -   all withdrawals from the bridge contract in this block
    -   Preconditions: block is final (e.g. 6 deep in the chain)
