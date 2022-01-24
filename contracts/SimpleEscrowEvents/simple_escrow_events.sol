pragma solidity ^0.4.0;
contract SimpleEscrowEvents {

    event Payment(string indexed contractAddress, uint amount);

    // send money to lock into escrow account
    function send(string contractAddress, uint amount) public {
        require(isStringsEqual(contractAddress, "0x0000000000000000000000000000000000000000") || isStringsEqual(contractAddress, "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), "contractAddress not supported");
        
        emit Payment(contractAddress, amount);
    }

    function isStringsEqual(string memory a, string memory b) public pure returns (bool) {
        return (keccak256(abi.encodePacked((a))) == keccak256(abi.encodePacked((b))));
    }
}