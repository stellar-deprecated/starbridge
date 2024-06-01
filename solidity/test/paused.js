const { expect } = require("chai");
const { ethers } = require("hardhat");
const { validTimestamp, expiredTimestamp } = require("./util");
const { updateSigners } = require("./updateSigners");

const PAUSE_NOTHING = 0;
const PAUSE_DEPOSITS = 1;
const PAUSE_WITHDRAWALS = 2;
const PAUSE_WITHDRAWALS_AND_DEPOSITS = 3;

let pausedNonce = 0;

function nextPauseNonce() {
    return pausedNonce++;
}

async function setPaused(bridge, signers, domainSeparator, state, nonce, expiration) {
    const request = [state, nonce, expiration];
    const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
        [domainSeparator, ethers.utils.id("setPaused"), request]
    )));
    const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
    return bridge.setPaused(request, signatures, [...Array(signers.length).keys()]);
}

describe("setPaused", function() {
    let signers;
    let bridge;
    let domainSeparator;

    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);
        domainSeparator = await bridge.domainSeparator();
    });

    it("is rejected if paused bitmask is invalid", async function() {
        await expect(
            setPaused(bridge, signers, domainSeparator, 4, nextPauseNonce(), validTimestamp())
        ).to.be.revertedWith("invalid paused value");
    });

    it("is rejected if expired", async function() {
        await expect(
            setPaused(bridge, signers, domainSeparator, PAUSE_DEPOSITS, nextPauseNonce(), expiredTimestamp())
        ).to.be.revertedWith("request is expired");
    });

    it("is rejected if method id is invalid", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused1"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.setPaused(request, signatures, [...Array(signers.length).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if domain separator is invalid", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version()) + 1, 31337, bridge.address] 
        ));
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [invalidDomainSeparator, ethers.utils.id("setPaused"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.setPaused(request, signatures, [...Array(signers.length).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if there are too few signatures", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,19).map(s => s.signMessage(hash)));
        await expect(
            bridge.setPaused(request, signatures, [...Array(signatures.length).keys()])
        ).to.be.revertedWith("not enough signatures");
    });

    it("is rejected if there are invalid signatures", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        signatures[0] = signatures[1];
        await expect(
            bridge.setPaused(request, signatures, [...Array(signatures.length).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("is rejected if the signatures are not sorted", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        const tmp = signatures[1];
        signatures[1] = signatures[0];
        signatures[0] = tmp; 
        const indexes = [...Array(20).keys()];
        indexes[0] = 1;
        indexes[1] = 0;
        await expect(
            bridge.setPaused(request, signatures, indexes)
        ).to.be.revertedWith("signatures not sorted by signer");
    });

    it("is rejected if the indexes length does not match signatures length", async function() {
        const request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.setPaused(request, signatures, [...Array(19).keys()])
        ).to.be.revertedWith("number of signatures does not equal number of indexes");
    });

    it("succeeds", async function() {
        let request = [PAUSE_DEPOSITS, 0, validTimestamp()];
        let hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        let signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await bridge.setPaused(request, signatures, [...Array(20).keys()]);

        // reusing transaction will be rejected
        await expect(
            bridge.setPaused(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("request is already fulfilled");
    });

    it("updateSigners invalidates setPaused transactions", async function() {
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);

        let request = [PAUSE_NOTHING, 0, validTimestamp()];
        let hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(uint8, uint256, uint256)"], 
            [domainSeparator, ethers.utils.id("setPaused"), request]
        )));
        let signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.setPaused(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });
});

module.exports = {
    PAUSE_NOTHING,
    PAUSE_DEPOSITS,
    PAUSE_WITHDRAWALS,
    PAUSE_WITHDRAWALS_AND_DEPOSITS,
    setPaused,
    nextPauseNonce,
};