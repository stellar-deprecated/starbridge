# Starbridge

This repository contains formal models related to the
[Starbridge](https://github.com/stellar/starbridge) project.

[starbridge.ivy](https://github.com/nano-o/Starbridge-modeling/blob/main/ivy/shared/starbridge.ivy)
contains a model of the Ethereum to Stellar deposit flow (with refunds) and
a safety proof in the form of an inductive invariant.

[starbridge-timelock.ivy](https://github.com/nano-o/Starbridge-modeling/blob/main/ivy/shared/starbridge-timelock.ivy)
contains a different model of the Ethereum to Stellar deposit flow based on an
idea of Tamir.

To check the proofs (substitute the file you want to check):

```
IVY_FILE=starbridge-timelock.ivy docker-compose run --project-directory ivy/ --rm starbridge-ivy
```

To plot the dependencies between invariants:

```
IVY_FILE=starbridge-timelock.ivy docker-compose run  --project-directory ivy/ --rm starbridge-ivy-poisonivy
```

You will then find a `png` file in `ivy/shared/`
