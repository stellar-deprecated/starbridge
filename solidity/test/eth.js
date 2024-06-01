const { expect } = require("chai");
const { ethers, waffle } = require("hardhat");
const { PAUSE_DEPOSITS, PAUSE_NOTHING, PAUSE_WITHDRAWALS_AND_DEPOSITS, setPaused, nextPauseNonce, PAUSE_WITHDRAWALS } = require("./paused");
const { updateSigners } = require("./updateSigners");
const { validTimestamp, expiredTimestamp } = require("./util");

describe("Deposit & Withdraw ETH", function() {
    let signers;
    let bridge;
    let domainSeparator;
    let sender;


    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        sender = signers[0];
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);
        domainSeparator = await bridge.domainSeparator();
    });

    it("fallback function reverts", async function() {
        await expect(signers[0].sendTransaction({to: bridge.address, value: ethers.utils.parseEther("1.0")})).to.be.reverted;
    });

    it("deposits of 0 ETH are rejected", async function() {
        await expect(bridge.depositETH(1, {value: 0})).to.be.revertedWith("deposit amount is zero");
    });

    it("deposits are rejected when bridge is paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_DEPOSITS, nextPauseNonce(), validTimestamp());
        await expect(bridge.depositETH(1, {value: ethers.utils.parseEther("1.0")})).to.be.revertedWith("deposits are paused");
        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    it("deposits is successful", async function() {
        const before = await waffle.provider.getBalance(bridge.address);
        await bridge.depositETH(1, {value: ethers.utils.parseEther("1.0")});
        const after = await waffle.provider.getBalance(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));  
    });

    it("deposits succeed when withdrawals are paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_WITHDRAWALS, nextPauseNonce(), validTimestamp());

        const before = await waffle.provider.getBalance(bridge.address);
        await bridge.depositETH(1, {value: ethers.utils.parseEther("1.0")});
        const after = await waffle.provider.getBalance(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));

        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    async function withdrawETH(id, domainSeparator, expiration, recipient, amount) {
        const request = [id, expiration, recipient, amount];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        return bridge.withdrawETH(request, signatures, [...Array(20).keys()]);
    }

    it("withdrawals are rejected when bridge is paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_WITHDRAWALS, nextPauseNonce(), validTimestamp());

        const recipient = signers[1].address;
        await expect(withdrawETH(
            ethers.utils.formatBytes32String("0"), 
            domainSeparator, 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("withdrawals are paused");

        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    it("withdrawals and deposits are rejected when bridge is paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_WITHDRAWALS_AND_DEPOSITS, nextPauseNonce(), validTimestamp());

        await expect(bridge.depositETH(1, {value: ethers.utils.parseEther("1.0")})).to.be.revertedWith("deposits are paused");

        const recipient = signers[1].address;
        await expect(withdrawETH(
            ethers.utils.formatBytes32String("0"), 
            domainSeparator, 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("withdrawals are paused");

        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    it("expired withdrawals are rejected", async function() {
        const recipient = signers[1].address;
        await expect(withdrawETH(
            ethers.utils.formatBytes32String("0"), 
            domainSeparator, 
            expiredTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("request is expired");
    });

    it("withdrawals with invalid method id are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH1"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawETH(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with invalid domain separator are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version()) + 1, 31337, bridge.address] 
        ));
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [invalidDomainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawETH(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with too few signatures are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,19).map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawETH(request, signatures, [...Array(signatures.length).keys()])
        ).to.be.revertedWith("not enough signatures");
    });

    it("withdrawals with invalid signatures are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        signatures[0] = signatures[1];
        await expect(
            bridge.withdrawETH(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with non sorted signatures are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        const tmp = signatures[1];
        signatures[1] = signatures[0];
        signatures[0] = tmp; 
        const indexes = [...Array(20).keys()];
        indexes[0] = 1;
        indexes[1] = 0;
        await expect(
            bridge.withdrawETH(request, signatures, indexes)
        ).to.be.revertedWith("signatures not sorted by signer");
    });

    it("withdrawals with indexes length not matching signatures length are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawETH"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawETH(request, signatures, [...Array(19).keys()])
        ).to.be.revertedWith("number of signatures does not equal number of indexes");
    });

    it("cannot withdraw more than bridge balance", async function() {
        await expect(
            withdrawETH(
                ethers.utils.formatBytes32String("0"),
                domainSeparator,
                validTimestamp(),
                signers[2].address,
                ethers.utils.parseEther("200")
            )
        ).to.be.revertedWith("ETH transfer failed");
    });

    it("withdrawal succeeds", async function() {
        let before = await waffle.provider.getBalance(bridge.address);
        await bridge.depositETH(1, {value: ethers.utils.parseEther("3.0")});
        let after = await waffle.provider.getBalance(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("3.0"));

        const recipient = signers[1].address;
        before = await waffle.provider.getBalance(recipient);
        await withdrawETH(ethers.utils.formatBytes32String("0"), domainSeparator, validTimestamp(), recipient, ethers.utils.parseEther("1.0"));
        after = await waffle.provider.getBalance(recipient);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));

        // reusing request id will be rejected
        await expect(
            withdrawETH(ethers.utils.formatBytes32String("0"), domainSeparator, validTimestamp(), signers[2].address, ethers.utils.parseEther("2.0"))
        ).to.be.revertedWith("request is already fulfilled");
    });

    it("updateSigners invalidates withdrawal transactions", async function() {
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);
        await expect(
            withdrawETH(ethers.utils.formatBytes32String("1"), domainSeparator, validTimestamp(), signers[2].address, ethers.utils.parseEther("1.0"))
        ).to.be.revertedWith("signature does not match");
    });
});