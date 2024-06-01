// contracts/Auth.sol
// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;
import "./StellarAsset.sol";


contract FeeToken is StellarAsset {

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_
    ) StellarAsset(name_, symbol_, decimals_) {}

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public virtual override returns (bool) {
        return super.transferFrom(sender, recipient, amount-1);
    }
}