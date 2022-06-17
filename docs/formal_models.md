[starbridge-timelock.ivy](https://github.com/nano-o/Starbridge-modeling/blob/main/ivy/shared/starbridge-timelock.ivy) contains a different model of the Ethereum to Stellar deposit flow based on an idea of Tamir.

To check the proofs:

```
IVY_FILE=starbridge-timelock.ivy docker-compose --project-directory docker/ run --rm starbridge-ivy
```

To plot the dependencies between invariants:

```
IVY_FILE=starbridge-timelock.ivy docker-compose --project-directory docker/ run --rm starbridge-ivy-poisonivy
```

You will then find a `png` file in `docker/output/`. It depicts a graph representing the dependencies between invariants (an invariant is inductive if all its children in the graph are assumed to hold). In some sense, it is a graphical representation of the proof.
