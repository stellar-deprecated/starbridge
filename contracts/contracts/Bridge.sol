// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;
import "./Auth.sol";
import "./StellarAsset.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

struct WithdrawERC20Request {
    bytes32 id;
    uint256 expiration;
    address recipient;
    address token;
    uint256 amount;
}

struct WithdrawETHRequest {
    bytes32 id;
    uint256 expiration;
    address recipient;
    uint256 amount;
}

struct CreateStellarAssetRequest {
    uint8 decimals;
    string name;
    string symbol;
}

contract Bridge is Auth {
    constructor(address[] memory _signers, uint8 _minThreshold) Auth(_signers, _minThreshold) {}

    event DepositERC20(
        address token,
        address sender,
        uint256 destination,
        uint256 amount
    );
    event DepositETH(
        address sender,
        uint256 destination,
        uint256 amount
    );

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

    event WithdrawERC20(
        bytes32 id,
        address recipient,
        address token,
        uint256 amount
    );
    event WithdrawETH(
        bytes32 id,
        address recipient,
        uint256 amount
    );

    event RegisterStellarAsset(address asset);

    mapping(address => bool) public isStellarAsset;

    function depositERC20(
        address token,
        uint256 destination,
        uint256 amount
    ) external {
        require(amount > 0);
        if (isStellarAsset[token]) {
            StellarAsset(token).burn(msg.sender, amount);
        } else {
            SafeERC20.safeTransferFrom(
                IERC20(token),
                msg.sender,
                address(this),
                amount
            );
        }
        emit DepositERC20(token, msg.sender, destination, amount);
    }

    function depositETH(uint256 destination) external payable {
        require(msg.value > 0);
        emit DepositETH(msg.sender, destination, msg.value);
    }

    function withdrawERC20(
        WithdrawERC20Request calldata request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        verifyRequest(
            keccak256(abi.encode(version, request)),
            request.id,
            request.expiration,
            signatures,
            indexes
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
        emit WithdrawERC20(
            request.id,
            request.recipient,
            request.token,
            request.amount
        );
    }

    function withdrawETH(
        WithdrawETHRequest calldata request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        verifyRequest(
            keccak256(abi.encode(version, request)),
            request.id, 
            request.expiration, 
            signatures,
            indexes);

        (bool success, ) = request.recipient.call{value: request.amount}("");
        require(success, "ETH transfer failed");
        emit WithdrawETH(request.id, request.recipient, request.amount);
    }

    function registerStellarAsset(
        CreateStellarAssetRequest memory request,
        bytes[] calldata signatures,
        uint8[] calldata indexes
    ) external {
        bytes32 requestHash = keccak256(abi.encode(
            version,
            request.decimals,
            keccak256(bytes(request.name)),
            keccak256(bytes(request.symbol))
        ));
        verifySignatures(
            requestHash,
            signatures,
            indexes
        );

        address asset = address(
            new StellarAsset{salt: requestHash}(
                request.name,
                request.symbol,
                request.decimals
            )
        );

        isStellarAsset[asset] = true;
        emit RegisterStellarAsset(asset);
    }
}
