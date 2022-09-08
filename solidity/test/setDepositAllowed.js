const { expect } = require("chai");
const { ethers } = require("hardhat");
const { expiredTimestamp, validTimestamp } = require("./util");
const { updateSigners } = require("./updateSigners");

async function setDepositAllowed(bridge, signers, domainSeparator, token, allowed, nonce, expiration) {
    const request = [token, allowed, nonce, expiration];
    const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
        [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
    )));
    const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
    return bridge.setDepositAllowed(request, signatures, [...Array(20).keys()]);
}

describe("setDepositAllowed", function() {
    let signers;
    let bridge;
    let domainSeparator;
    const ethAddress = "0x0000000000000000000000000000000000000000";

    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);
        domainSeparator = await bridge.domainSeparator();
    });

    it("is rejected if expired", async function() {
        await expect(
            setDepositAllowed(bridge, signers, domainSeparator, ethAddress, false, 0, expiredTimestamp())
        ).to.be.revertedWith("request is expired");
    });

    it("is rejected if method id is invalid", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed1"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if domain separator is invalid", async function() {
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version())+1, 31337, bridge.address] 
        ));
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [invalidDomainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if there are too few signatures", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,19).map(s => s.signMessage(hash)));
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(signatures.length).keys()])
        ).to.be.revertedWith("not enough signatures");
    });

    it("is rejected if there are invalid signatures", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        signatures[0] = signatures[1];
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(signatures.length).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if the signatures are not sorted", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        const tmp = signatures[1];
        signatures[1] = signatures[0];
        signatures[0] = tmp; 
        const indexes = [...Array(20).keys()];
        indexes[0] = 1;
        indexes[1] = 0;
        await expect(
            bridge.setDepositAllowed(request, signatures, indexes)
        ).to.be.revertedWith("signatures not sorted by signer");
    });

    it("is rejected if the indexes length does not match signatures length", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(19).keys()])
        ).to.be.revertedWith("number of signatures does not equal number of indexes");
    });

    it("nonce prevents transaction replay", async function() {
        const request = [ethAddress, false, 0, validTimestamp()];
        let hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        let signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await bridge.setDepositAllowed(request, signatures, [...Array(20).keys()]);

        // reusing transaction will be rejected
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("request is already fulfilled");
    });

    it("updateSigners invalidates setDepositAllowed transactions", async function() {
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);

        const request = [ethAddress, true, 1, validTimestamp()];
        let hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(address, bool, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setDepositAllowed"), request]
        )));
        let signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.setDepositAllowed(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });
});

module.exports = {
    setDepositAllowed,
};