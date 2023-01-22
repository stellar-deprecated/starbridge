package soroban_bridge

import _ "embed"

//go:generate sh -c "cargo build --target wasm32-unknown-unknown --release"

//go:embed target/wasm32-unknown-unknown/release/soroban_bridge.wasm
var SorobanBridgeWasm []byte
