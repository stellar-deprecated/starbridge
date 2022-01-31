// contracts/Auth.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Auth.sol";
import "./StellarAsset.sol";

contract StellarAssetFactory {
    mapping(bytes32 => address) public nameHashToAsset;
    mapping(address => bool) public isStellarAsset;

    event CreateStellarAsset(address asset);

    function createStellarAsset(
        MintStellarAssetRequest memory request,
        bytes32 nameHash
    ) internal returns (address asset) {
        asset = address(
            new StellarAsset{salt: nameHash}(
                request.name,
                request.symbol,
                request.decimals
            )
        );
        nameHashToAsset[nameHash] = asset;
        isStellarAsset[asset] = true;
        emit CreateStellarAsset(asset);
    }
}
