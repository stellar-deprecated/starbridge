// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract Auth {
    address[] public signers;
    uint8 public minThreshold;
     // every time the validator set is updated, the version is incremented.
     // increasing version numbers ensure that updateSigners() signatures cannot be reused
    uint256 public nextVersion;
    event RegisterSigners(uint256 version, address[] signers, uint8 minThreshold);

    mapping(bytes32 => bool) private fulfilledrequests;


    constructor(address[] memory _signers, uint8 _minThreshold) {
        _updateSigners(0, _signers, _minThreshold);
        nextVersion = 1;
    }

    function _updateSigners(uint256 curVersion, address[] memory _signers, uint8 _minThreshold) internal {
        require(_signers.length > 0, "too few signers");
        require(_signers.length < 256, "too many signers");
        require(_minThreshold > _signers.length / 2, "min threshold is too low");
        require(_minThreshold <= _signers.length, "min threshold is too high");
        // by requiring signers to be sorted we can verify there are no duplicate
        // signers in linear time
        for (uint8 i = 1; i < _signers.length; i++) {
            require(_signers[i-1] < _signers[i], "signers not sorted");
        }
        
        signers = _signers;
        minThreshold = _minThreshold;
        emit RegisterSigners(curVersion, _signers, _minThreshold);
    }

    function updateSigners(
        address[] calldata _signers, 
        uint8 _minThreshold,
        bytes[] calldata signatures, 
        uint8[] calldata indexes
    ) external {
        uint256 curVersion = nextVersion++;
        bytes32 h = keccak256(abi.encode(curVersion, _signers, _minThreshold));
        verifySignatures(h, signatures, indexes);
        _updateSigners(curVersion, _signers, _minThreshold);
    }

    function verifySignatures(bytes32 h, bytes[] memory signatures, uint8[] memory indexes)
        internal view
    {
        require(
            signatures.length == indexes.length,
            "number of signatures does not equal number of indexes"
        );
        require(signatures.length >= minThreshold, "not enough signatures");
        uint8 prev = 0;
        for (uint8 i = 0; i < signatures.length; i++) {
            uint8 idx = indexes[i];
            address signer = signers[idx];
            // by requiring indexes to be sorted we can verify there are no duplicate
            // signatures in linear time
            require(i == 0 || idx > prev, "signatures not sorted by signer");
            require(
                ECDSA.recover(h, signatures[i]) == signer,
                "signature does not match"
            );
            prev = idx;
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
        require(block.timestamp < expiration, "request is expired");
    }

    function requestStatus(bytes32 requestID) external view returns (bool, uint256) {
        return (fulfilledrequests[requestID], block.number);
    }
}
