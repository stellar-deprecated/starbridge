pragma solidity ^0.4.0;
contract SimpleEscrowEvents {

    event Payment(string indexed destinationStellarAddress, string indexed tokenContractAddress, uint tokenAmount);

    // send money to lock into escrow account
    function send(string destinationStellarAddress, string tokenContractAddress, uint tokenAmount) public {
        require(isStringsEqual(tokenContractAddress, "0x0000000000000000000000000000000000000000") || isStringsEqual(tokenContractAddress, "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), "tokenContractAddress not supported");
        
        emit Payment(destinationStellarAddress, tokenContractAddress, tokenAmount);
    }

    function isStringsEqual(string memory a, string memory b) private pure returns (bool) {
        return (keccak256(abi.encodePacked((a))) == keccak256(abi.encodePacked((b))));
    }
}