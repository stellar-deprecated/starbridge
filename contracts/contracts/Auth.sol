// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract Auth {
    address[] public signers;
    uint8 minThreshold;
     // every time the validator set is updated, the version is incremented.
     // increasing version numbers ensure that updateSigners() signatures cannot be reused
    uint256 nextVersion;
    event RegisterSigners(uint256 version, address[] signers, uint8 minThreshold);

    mapping(bytes32 => bool) public fulfilledrequests;


    constructor(address[] memory _signers, uint8 _minThreshold) {
        _updateSigners(_signers, _minThreshold);
    }

    function _updateSigners(address[] memory _signers, uint8 _minThreshold) internal {
        require(_signers.length < 256, "too many signers");
        require(_minThreshold > signers.length / 2, "min threshold is too low");
        // by requiring signers to be sorted we can verify there are no duplicate
        // signers in linear time
        for (uint8 i = 1; i < _signers.length; i++) {
            require(_signers[i-1] < _signers[i], "signers not sorted");
        }
        
        signers = _signers;
        minThreshold = _minThreshold;
        emit RegisterSigners(nextVersion++, _signers, _minThreshold);
    }

    function updateSigners(
        address[] calldata _signers, 
        uint8 _minThreshold,
        bytes[] calldata signatures, 
        uint8[] calldata indexes
    ) external {
        bytes32 h = keccak256(abi.encode(nextVersion, _signers, _minThreshold));
        verifySignatures(h, signatures, indexes);
        _updateSigners(_signers, _minThreshold);
    }

    function verifySignatures(bytes32 h, bytes[] memory signatures, uint8[] memory indexes)
        internal view
    {
        require(
            signatures.length == indexes.length,
            "number of signatures does not equal number of signers"
        );
        require(signatures.length >= minThreshold, "not enough signatures");
        address prev;
        for (uint256 i = 0; i < signatures.length; i++) {
            address signer = signers[indexes[i]];
            // by requiring indexes to be sorted we can verify there are no duplicate
            // signatures in linear time
            require(signer > prev, "signatures not sorted by signer");
            require(
                ECDSA.recover(h, signatures[i]) == signer,
                "signature does not match"
            );
            prev = signer;
        }
    }

    function verifyRequest(
        bytes32 requestHash,
        bytes32 requestID,
        uint256 expiration,
        bytes[] memory signatures,
        uint8[] memory indexes
    ) internal {
        verifySignatures(requestHash, signatures, indexes);
        require(!fulfilledrequests[requestID], "request is already fulfilled");
        fulfilledrequests[requestID] = true;
        require(block.number < expiration, "request is expired");
    }
}
