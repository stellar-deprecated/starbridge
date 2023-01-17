# Setup

See https://soroban.stellar.org/docs/getting-started/setup for setting up your dev environment.

# Building

Run this command from the `soroban-bridge` directory:

```
cargo build --target wasm32-unknown-unknown --release
ls target/wasm32-unknown-unknown/release/soroban_bridge.wasm
```

The byte code of the `soroban-bridge` smart contract will be in `target/wasm32-unknown-unknown/release/soroban_bridge.wasm`.