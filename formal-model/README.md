# Starbridge

This repository contains a formal model of the protocol to transfer native
Ethereum tokens to Stellar (including refunds) and a safety proof in the form
of an inductive invariant that is checked by
[Ivy](https://github.com/kenmcmil/ivy).

To check the proof, run `make check-proof` (tested on Ubuntu; you must have
docker installed):

To plot the dependencies between invariants, which gives a kind of visual
representation of the proof, run `make dependency-graph`. You will then find a
`png` file in `ouput/` depicting a directed graph whose nodes are invariants
(identified by their line number in `starbridge.ivy`). The meaning of the graph
is that an invariant is inductive relatively to its children in the graph (i.e.
if one assumes `I & C`, where `C` is the conjuntion of the children of `I`, and
on takes a protocol step, then `I` holds again).
