package solidity

//go:generate sh -c "cd ../solidity && npm clean-install && npm exec -- hardhat compile"

//go:generate sh -c "jq '.abi' ../solidity/artifacts/contracts/Auth.sol/Auth.json | go run github.com/ethereum/go-ethereum/cmd/abigen@v1.10.15 --abi - --pkg solidityauth --type Auth --out ./solidityauth/auth.go"
//go:generate sh -c "jq '.abi' ../solidity/artifacts/contracts/Bridge.sol/Bridge.json | go run github.com/ethereum/go-ethereum/cmd/abigen@v1.10.15 --abi - --pkg soliditybridge --type Bridge --out ./soliditybridge/bridge.go"
//go:generate sh -c "jq '.abi' ../solidity/artifacts/contracts/StellarAsset.sol/StellarAsset.json | go run github.com/ethereum/go-ethereum/cmd/abigen@v1.10.15 --abi - --pkg soliditystellarasset --type StellarAsset --out ./soliditystellarasset/stellar_asset.go"
//go:generate sh -c "jq '.abi' ../solidity/artifacts/contracts/StellarAssetFactory.sol/StellarAssetFactory.json | go run github.com/ethereum/go-ethereum/cmd/abigen@v1.10.15 --abi - --pkg soliditystellarassetfactory --type StellarAssetFactory --out ./soliditystellarassetfactory/stellar_asset_factory.go"
