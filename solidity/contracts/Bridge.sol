// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;
import "./Auth.sol";
import "./StellarAsset.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";

// Every bridge transfer has a globally unique id.
//
// For a ETH -> Stellar transfer the id would be keccak256(abi.encode(txHash, logIndex))
// where txHash is the hash of the ethereum transaction containing the deposit event and
// logIndex is the index of the deposit event in the ethereum block.
// The logIndex is necessary because a single ethereum transaction could call
// depositERC20() / depositETH() multiple times.
// The bridge smart contract will ensure that a ETH -> Stellar transfer can only be refunded
// once by maintaining a set of fulfilled ids.
//
// For a Stellar -> ETH transfer the id wouold be Stellar transaction hash assuming a
// stellar transaction can only contain one Stellar -> ETH transfer.
// Similar to the refund case, the bridge smart contract will ensure that a Stellar -> ETH
// transfer can only be completed once by mantaining a set of fulfilled ids.

// WithdrawERC20Request is the payload for the withdrawERC20() transaction.
struct WithdrawERC20Request {
    bytes32 id; // the id of the transfer
    uint256 expiration; // unix timestamp of when the transaction should expire
    address recipient; // ethereum address who will receive the ERC20 tokens
    address token; // the ERC20 token
    uint256 amount; // the amount of tokens to be transferred
}

// WithdrawETHRequest is the payload for the withdrawETH() transaction.
struct WithdrawETHRequest {
    bytes32 id; // the id of the transfer
    uint256 expiration; // unix timestamp of when the transaction should expire
    address recipient; // ethereum address who will receive the ETH
    uint256 amount; // the amount of ETH to be transferred
}

// SetPausedRequest is the payload for the setPaused() transaction.
struct SetPausedRequest {
    uint8 value; // bitmask representing whether deposits / withdrawals are enabled
    uint256 nonce; // used to make each transaction unique for replay prevention
    uint256 expiration; // unix timestamp of when the transaction should expire
}

// RegisterStellarAssetRequest is the payload for the registerStellarAsset() transaction.
// The three fields define a new ERC20 token which represents the ethereum equivalent of
// a Stellar asset.
struct RegisterStellarAssetRequest {
    uint8 decimals;
    string name;
    string symbol;
}

// bitwise flag representing the state where no deposits are allowed on the bridge
uint8 constant PAUSE_DEPOSITS = 1 << 0;
// bitwise flag representing the state where no withdrawals are allowed on the bridge
uint8 constant PAUSE_WITHDRAWALS =  1 << 1;
// bitwise flag representing the state where no withdrawals or deposits are allowed on the bridge
uint8 constant PAUSE_DEPOSITS_AND_WITHDRAWALS = PAUSE_DEPOSITS | PAUSE_WITHDRAWALS;

// SET_PAUSED_ID is used to distinguish setPaused() signatures from signatures for other bridge functions.
bytes32 constant SET_PAUSED_ID = keccak256("setPaused");
// REGISTER_STELLAR_ASSET_ID is used to distinguish registerStellarAsset() signatures from signatures for other bridge functions.
bytes32 constant REGISTER_STELLAR_ASSET_ID = keccak256("registerStellarAsset");
// WITHDRAW_ETH_ID is used to distinguish withdrawETH() signatures from signatures for other bridge functions.
bytes32 constant WITHDRAW_ETH_ID = keccak256("withdrawETH");
// WITHDRAW_ERC20_ID is used to distinguish withdrawERC20() signatures from signatures for other bridge functions.
bytes32 constant WITHDRAW_ERC20_ID = keccak256("withdrawERC20");

contract Bridge is Auth, ReentrancyGuard {
    // paused is a bitmask which determines whether deposits / withdrawals are enabled on the bridge
    uint8 public paused;
    // SetPaused is emitted whenever the paused state of the bridge changes
    event SetPaused(uint8 value);

    // to create a Bridge instance you need to provide the validator set configuration
    constructor(address[] memory _signers, uint8 _minThreshold) Auth(_signers, _minThreshold) {}

    // Deposit is emitted whenever ERC20 tokens (or ETH) are deposited on the bridge.
    // The Deposit event initiates a Ethereum -> Stellar transfer.
    event Deposit(
        // 0x0 coresponds to ETH
        address token,
        address sender,
        uint256 destination,
        uint256 amount
    );

    // Withdraw is emitted whenever ERC20 tokens (or ETH) is claimed from the bridge.
    // The Withdraw event corresponds to completing a Stellar -> Ethereum transfer or
    // refunding a Ethereum -> Stellar transfer.
    event Withdraw(
        bytes32 id,
        // 0x0 coresponds to ETH
        address token,
        address recipient,
        uint256 amount
    );

    // RegisterStellarAsset is emitted whenever an ERC20 token is created
    // to represent a Stellar asset.
    event RegisterStellarAsset(address asset);

    // isStellarAsset identifies whether an ERC20 token is a Stellar asset
    // created by the bridge.
    mapping(address => bool) public isStellarAsset;

    // depositERC20() deposits ERC20 tokens to the bridge and starts a ERC20 -> Stellar
    // transfer. If deposits are disabled this function will fail.
    function depositERC20(
        address token,
        uint256 destination,
        uint256 amount
    ) external nonReentrant {
        require((paused & PAUSE_DEPOSITS) == 0, "deposits are paused");
        require(amount > 0, "deposit amount is zero");

        emit Deposit(token, msg.sender, destination, amount);

        if (isStellarAsset[token]) {
            StellarAsset(token).burn(msg.sender, amount);
        } else {
            IERC20 tokenContract = IERC20(token);
            uint256 before = tokenContract.balanceOf(address(this));
            SafeERC20.safeTransferFrom(
                tokenContract,
                msg.sender,
                address(this),
                amount
            );
            uint256 received = SafeMath.sub(tokenContract.balanceOf(address(this)), before);
            require(amount == received, "received amount not equal to expected amount");
        }
    }

    // depositETH() deposits ETH to the bridge and starts a ETH -> Stellar
    // transfer. If deposits are disabled this function will fail.
    function depositETH(uint256 destination) external payable {
        require((paused & PAUSE_DEPOSITS) == 0, "deposits are paused");
        require(msg.value > 0, "deposit amount is zero");
        emit Deposit(address(0), msg.sender, destination, msg.value);
    }

    // withdrawERC20() claims ERC20 tokens from the bridge. This can correspond to
    // fulfilling a Stellar -> ERC20 transfer or refunding a ERC20 -> Stellar transfer.
    // If withdrawals are disabled this function will fail.
    // withdrawERC20() must be authorized by the bridge validators otherwise the transaction
    // will fail. Replay prevention is implemented using the transfer id.
    // Any attempts to withdraw multiple times for the same transfer will fail.
    function withdrawERC20(
        WithdrawERC20Request calldata request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        require((paused & PAUSE_WITHDRAWALS) == 0, "withdrawals are paused");
        verifyRequest(
            keccak256(abi.encode(domainSeparator, WITHDRAW_ERC20_ID, request)),
            request.id,
            request.expiration,
            signatures,
            indexes
        );
        emit Withdraw(
            request.id,
            request.token,
            request.recipient,
            request.amount
        );

        if (isStellarAsset[request.token]) {
            StellarAsset(request.token).mint(request.recipient, request.amount);
        } else {
            SafeERC20.safeTransfer(
                IERC20(request.token),
                request.recipient,
                request.amount
            );
        }
    }

    // withdrawETH() claims ETH tokens from the bridge. This can correspond to
    // fulfilling a Stellar -> ETH transfer or refunding a ETH -> Stellar transfer.
    // If withdrawals are disabled this function will fail.
    // withdrawETH() must be authorized by the bridge validators otherwise the transaction
    // will fail. Replay prevention is implemented using the transfer id.
    // Any attempts to withdraw multiple times for the same transfer will fail.
    function withdrawETH(
        WithdrawETHRequest calldata request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        require((paused & PAUSE_WITHDRAWALS) == 0, "withdrawals are paused");
        verifyRequest(
            keccak256(abi.encode(domainSeparator, WITHDRAW_ETH_ID, request)),
            request.id, 
            request.expiration, 
            signatures,
            indexes);

        emit Withdraw(request.id, address(0), request.recipient, request.amount);
        (bool success, ) = request.recipient.call{value: request.amount}("");
        require(success, "ETH transfer failed");
    }

    // setPaused() will enable or disable withdrawals / deposits.
    // setPaused() must be authorized by the bridge validators otherwise the transaction
    // will fail. Replay prevention is implemented by storing the request hash in the
    // fulfilledrequests set.
    function setPaused(
        SetPausedRequest memory request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        require(request.value <= PAUSE_DEPOSITS_AND_WITHDRAWALS, "invalid paused value");
        bytes32 requestHash = keccak256(abi.encode(domainSeparator, SET_PAUSED_ID, request));
        // ensure the same setPaused() transaction cannot be used more than once
        verifyRequest(requestHash, requestHash, request.expiration, signatures, indexes);
        emit SetPaused(request.value);
        paused = request.value;
    }

    // registerStellarAsset() will creates an ERC20 token to represent a stellar asset.
    // registerStellarAsset() must be authorized by the bridge validators otherwise the transaction
    // will fail. Replay prevention is impemented by creating the ERC20 via the CREATE2 opcode (see
    // https://eips.ethereum.org/EIPS/eip-1014 )
    function registerStellarAsset(
        RegisterStellarAssetRequest memory request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        bytes32 requestHash = keccak256(abi.encode(
            domainSeparator,
            REGISTER_STELLAR_ASSET_ID,
            request.decimals,
            keccak256(bytes(request.name)),
            keccak256(bytes(request.symbol))
        ));
        verifySignatures(
            requestHash,
            signatures,
            indexes
        );

        bytes32 assetHash = keccak256(abi.encode(
            request.decimals,
            keccak256(bytes(request.name)),
            keccak256(bytes(request.symbol))
        ));
        address asset = address(
            new StellarAsset{salt: assetHash}(
                request.name,
                request.symbol,
                request.decimals
            )
        );

        emit RegisterStellarAsset(asset);
        isStellarAsset[asset] = true;
    }
}
