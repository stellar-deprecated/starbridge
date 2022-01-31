// contracts/Auth.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Auth.sol";
import "./StellarAsset.sol";
import "./StellarAssetFactory.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract Bridge is Auth, StellarAssetFactory {
    constructor(address[] memory _signers) Auth(_signers) {}

    event DepositERC20(
        address token,
        address sender,
        uint256 destination,
        uint256 amount
    );
    event DepositETH(address sender, uint256 destination, uint256 amount);
    event BurnStellarAsset(address asset, address holder, uint256 amount);

    event WithdrawERC20(
        uint256 nonce,
        address recipient,
        address token,
        uint256 amount
    );
    event WithdrawETH(uint256 nonce, address recipient, uint256 amount);
    event MintStellarAsset(
        uint256 nonce,
        address asset,
        address recipient,
        uint256 amount
    );

    function depositERC20(
        address token,
        uint256 destination,
        uint256 amount
    ) external {
        if (isStellarAsset[token]) {
            StellarAsset(token).burn(msg.sender, amount);
            emit BurnStellarAsset(token, msg.sender, amount);
        } else {
            SafeERC20.safeTransferFrom(
                IERC20(token),
                msg.sender,
                address(this),
                amount
            );
            emit DepositERC20(token, msg.sender, destination, amount);
        }
    }

    function depositETH(uint256 destination) external payable {
        require(msg.value > 0);
        emit DepositETH(msg.sender, destination, msg.value);
    }

    function withdrawERC20(
        WithdrawERC20Request calldata request,
        bytes[] calldata signatures
    ) external {
        verifyRequest(hashWithdrawERC20Request(request), signatures);
        SafeERC20.safeTransfer(
            IERC20(request.token),
            request.recipient,
            request.amount
        );
        emit WithdrawERC20(
            request.nonce,
            request.recipient,
            request.token,
            request.amount
        );
    }

    function withdrawETH(
        WithdrawETHRequest calldata request,
        bytes[] calldata signatures
    ) external {
        verifyRequest(hashWithdrawETHRequest(request), signatures);
        (bool success, ) = request.recipient.call{value: request.amount}("");
        require(success, "ETH transfer failed");
        emit WithdrawETH(request.nonce, request.recipient, request.amount);
    }

    function mintStellarAsset(
        MintStellarAssetRequest calldata request,
        bytes[] calldata signatures
    ) external {
        (
            InternalMintStellarAssetRequest memory internalRequest,
            bytes32 nameHash,
            bytes32 symbolHash
        ) = toInternalMintStellarAssetRequest(request);
        verifyRequest(
            hashMintStellarAssetRequest(internalRequest, nameHash, symbolHash),
            signatures
        );
        address stellarAsset = nameHashToAsset[nameHash];
        if (stellarAsset == address(0)) {
            stellarAsset = createStellarAsset(request, nameHash);
        } else {
            require(
                keccak256(
                    bytes(StellarAsset(stellarAsset).symbol())
                ) == keccak256(bytes((request.symbol))),
                "symbol does not match"
            );
            require(
                StellarAsset(stellarAsset).decimals() == request.decimals,
                "decimals does not match"
            );
        }
        StellarAsset(stellarAsset).mint(request.recipient, request.amount);
        emit MintStellarAsset(
            request.nonce,
            stellarAsset,
            request.recipient,
            request.amount
        );
    }
}
