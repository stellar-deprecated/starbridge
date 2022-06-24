package solidity

// Generating go code for the solidity contracts requires that the following
// tools be installed:
//   - npm
//   - jq
//   - go

//go:generate sh -c "cd ../solidity && npm ci && npx hardhat compile"

//go:generate sh -c "jq '.abi' ../solidity/artifacts/contracts/Bridge.sol/Bridge.json | go run github.com/ethereum/go-ethereum/cmd/abigen@v1.10.15 --abi - --pkg solidity --type Bridge --out ./bridge.go"
