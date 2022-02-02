// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/draft-EIP712.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

struct WithdrawERC20Request {
    uint256 nonce; // stellar sequence number, op index of tx
    address recipient;
    address token;
    uint256 amount;
}
struct WithdrawETHRequest {
    uint256 nonce; // stellar sequence number, op index of tx
    address recipient;
    uint256 amount;
}
struct MintStellarAssetRequest {
    uint256 nonce; // stellar sequence number, op index of tx
    address recipient;
    uint256 amount;
    uint8 decimals;
    string name; // must be unique among all stellar assets
    string symbol;
}

contract Auth is EIP712("Stellar Bridge", "1") {
    address[] public signers;
    mapping(bytes32 => bool) private _usedHashes;

    event RegisterSigners(address[] signers);

    constructor(address[] memory _signers) {
        // TODO verify signers are unique
        signers = _signers;
        emit RegisterSigners(_signers);
    }

    function DOMAIN_SEPARATOR() external view returns (bytes32) {
        return _domainSeparatorV4();
    }

    bytes32 public constant WITHDRAW_ERC20_TYPEHASH =
        keccak256(
            "WithdrawERC20Request(uint256 nonce,address recipient,address token,uint256 amount)"
        );

    bytes32 public constant WITHDRAW_ETH_TYPEHASH =
        keccak256(
            "WithdrawETHRequest(uint256 nonce,address recipient,uint256 amount)"
        );

    struct InternalMintStellarAssetRequest {
        uint256 nonce; // stellar sequence number, op index of tx
        address recipient;
        uint256 amount;
        uint8 decimals;
    }

    bytes32 public constant MINT_STELLAR_ASSET_TYPEHASH =
        keccak256(
            "MintStellarAssetRequest(uint256 nonce,address recipient,uint256 amount,uint8 decimals)"
        );

    function toInternalMintStellarAssetRequest(
        MintStellarAssetRequest memory request
    )
        internal
        pure
        returns (
            InternalMintStellarAssetRequest memory,
            bytes32,
            bytes32
        )
    {
        InternalMintStellarAssetRequest memory internalRequest;
        assembly {
            internalRequest := request
        }

        return (
            internalRequest,
            keccak256(bytes(request.name)),
            keccak256(bytes(request.symbol))
        );
    }

    function hashWithdrawERC20Request(WithdrawERC20Request memory request)
        public
        view
        returns (bytes32)
    {
        return
            _hashTypedDataV4(
                keccak256(abi.encode(WITHDRAW_ERC20_TYPEHASH, request))
            );
    }

    function hashWithdrawETHRequest(WithdrawETHRequest memory request)
        public
        view
        returns (bytes32)
    {
        return
            _hashTypedDataV4(
                keccak256(abi.encode(WITHDRAW_ETH_TYPEHASH, request))
            );
    }

    function hashMintStellarAssetRequest(
        InternalMintStellarAssetRequest memory request,
        bytes32 nameHash,
        bytes32 symbolHash
    ) public view returns (bytes32) {
        return
            _hashTypedDataV4(
                keccak256(
                    abi.encode(
                        MINT_STELLAR_ASSET_TYPEHASH,
                        request,
                        nameHash,
                        symbolHash
                    )
                )
            );
    }

    function verifyRequest(bytes32 requestHash, bytes[] memory signatures)
        internal
    {
        require(
            signatures.length == signers.length,
            "number of signatures does not equal number of signers"
        );
        for (uint256 i = 0; i < signatures.length; i++) {
            require(
                ECDSA.recover(requestHash, signatures[i]) == signers[i],
                "signature does not match"
            );
        }
        require(
            !_usedHashes[requestHash],
            "request hash has already been used"
        );
        _usedHashes[requestHash] = true;
    }
}
