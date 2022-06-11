// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

// UPDATE_SIGNERS_ID is used to distinguish updateSigners() signatures from signatures for other bridge functions.
bytes32 constant UPDATE_SIGNERS_ID = keccak256("updateSigners");

contract Auth {
    // The validator set is configured by the following three fields:
    // signers - the full list of signers who can approve a bridge transaction.
    // Each of the signers has equal wieght.
    // minThreshold - the minimum amount of signers who need to approve a bridge transaction
    // for it to be valid.
    // version - a sequence number associated with the validator set. Whenver the validator
    // set configuration is updated the version will increment. The current version is part of the
    // payload for each bridge transaction. So, whenever the version is bumped that will invalidate
    // any previously signed bridge transactions.
    address[] public signers;
    uint8 public minThreshold;
    uint256 public version;
    // RegisterSigners is emitted whenever the validator set configuration is modified.
    event RegisterSigners(uint256 version, address[] signers, uint8 minThreshold);
    // fulfilledrequests is a set of all bridge requests which have been completed. This
    // set is used to prevent an attacker from replaying bridge transactions.
    mapping(bytes32 => bool) private fulfilledrequests;


    constructor(address[] memory _signers, uint8 _minThreshold) {
        _updateSigners(0, _signers, _minThreshold);
    }

    function _updateSigners(uint256 newVersion, address[] memory _signers, uint8 _minThreshold) internal {
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
        emit RegisterSigners(newVersion, _signers, _minThreshold);
    }

    // updateSigners() is called to update the validator set configuration for the bridge.
    // updateSigners() will bump the version field as a side effect and will emit a RegisterSigners
    // event detailing the new configuration.
    // The updateSigners() transactions must be authorized by the previous validator set.
    // The transactions cannot be replayed because the version field is incremented in updateSigners()
    // which invalidates all transaction which sign a payload containing older versions.
    function updateSigners(
        address[] calldata _signers, 
        uint8 _minThreshold,
        bytes[] calldata signatures, 
        uint8[] calldata indexes
    ) external {
        uint256 newVersion = ++version;
        bytes32 h = keccak256(abi.encode(newVersion-1, UPDATE_SIGNERS_ID, _signers, _minThreshold));
        verifySignatures(h, signatures, indexes);
        _updateSigners(newVersion, _signers, _minThreshold);
    }

    // verifySignatures() ensure that provided list of signatures map to the validator set
    // configured for the bridge and that there are at least minThreshold signers.
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

    // verifyRequest() will ensure the following three invariants hold
    // 1) the request is authorized by the bridge validators and a sufficient
    // number of signatures from the bridge validators.
    // 2) the request has not been executed before (replay protection).
    // 3) the request is not expired
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

    // requestStatus() takes a request id and returns whether it was fulfilled
    // along with the current block.
    // This function will be invoked by the validators using https://eth.wiki/json-rpc/API#eth_call
    // in order to check whether a bridge withdrawal / refund was executed.
    // Returning the current block is also necessary so that the validator
    // can associate the requestStatus() response with a moment in time.
    // For example, consider a Stellar -> ETH bridge deposit which occurred at time t.
    // Assume the validator calls requestStatus() on the transfer id and the response is (false, n).
    // If the timestamp of block n is greater than t + 24h then we can infer that the transfer
    // was not claimed on the ETH side of the bridge and that it is safe to authorize a refund on
    // the Stellar side.
    function requestStatus(bytes32 requestID) external view returns (bool, uint256) {
        return (fulfilledrequests[requestID], block.number);
    }
}
