# Starbridge

This repository contains a formal model of the protocol to transfer native
Ethereum tokens to Stellar (including refunds) and a safety proof in the form
of an inductive invariant that is checked by
[Ivy](https://github.com/kenmcmil/ivy).

To check the proof, run `make check-proof` (tested on Ubuntu; you must have
docker installed):

To plot the dependencies between invariants, run `make dependency-graph`. You
will then find a `png` file in `ouput/`
