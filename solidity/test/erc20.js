const { expect } = require("chai");
const { ethers } = require("hardhat");
const { PAUSE_DEPOSITS, PAUSE_NOTHING, PAUSE_WITHDRAWALS_AND_DEPOSITS, setPaused, nextPauseNonce, PAUSE_WITHDRAWALS } = require("./paused");
const { updateSigners } = require("./updateSigners");
const { validTimestamp, expiredTimestamp } = require("./util");
const { setDepositAllowed } = require("./setDepositAllowed");

async function withdrawERC20(bridge, token, signers, id, domainSeparator, expiration, recipient, amount) {
    const request = [id, expiration, recipient, token.address, amount];
    const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
        [domainSeparator, ethers.utils.id("withdrawERC20"), request]
    )));
    const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
    return bridge.withdrawERC20(request, signatures, [...Array(20).keys()]);
}

describe("Deposit & Withdraw ERC20", function() {
    let signers;
    let bridge;
    let token;
    let domainSeparator;
    let ERC20;
    let sender;

    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        sender = signers[0];
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);
        domainSeparator = await bridge.domainSeparator();

        ERC20 = await ethers.getContractFactory("StellarAsset");
        token = await ERC20.deploy("Test Token", "TEST", 18);
        await token.mint(sender.address, ethers.utils.parseEther("100.0"));
        await token.approve(bridge.address, ethers.utils.parseEther("300.0"));
        await setDepositAllowed(bridge, signers, domainSeparator, token.address, true, 0, validTimestamp());
    });

    it("deposits of 0 are rejected", async function() {
        await expect(bridge.depositERC20(token.address, 1, 0)).to.be.revertedWith("deposit amount is zero");
    });

    it("deposits are rejected when bridge is paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_DEPOSITS, nextPauseNonce(), validTimestamp());
        
        await expect(bridge.depositERC20(
            token.address, 1, ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("deposits are paused");

        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    it("block deposits for a specific ERC20 token", async function() {
        const blockedToken = await ERC20.deploy("Blocked Test Token", "BLOCKED", 18);
        await blockedToken.mint(sender.address, ethers.utils.parseEther("100.0"));
        await blockedToken.approve(bridge.address, ethers.utils.parseEther("300.0"));

        expect(await bridge.depositAllowed(blockedToken.address)).to.be.false;
        await expect(bridge.depositERC20(
            blockedToken.address, 1, ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("deposits not allowed for token");

        await setDepositAllowed(bridge, signers, domainSeparator, blockedToken.address, true, 1, validTimestamp());
        expect(await bridge.depositAllowed(blockedToken.address)).to.be.true;

        const before = await token.balanceOf(bridge.address);

        await bridge.depositERC20(
            blockedToken.address, 1, ethers.utils.parseEther("1.0")
        );

        const after = await blockedToken.balanceOf(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));

        await setDepositAllowed(bridge, signers, domainSeparator, blockedToken.address, false, 2, validTimestamp());

        expect(await bridge.depositAllowed(blockedToken.address)).to.be.false;
        await expect(bridge.depositERC20(
            blockedToken.address, 1, ethers.utils.parseEther("1.0")
        )).to.be.revertedWith("deposits not allowed for token");
    });

    it("cannot deposit more tokens than current balance", async function() {
        await expect(bridge.depositERC20(
            token.address, 1, ethers.utils.parseEther("200")
        )).revertedWith("ERC20: transfer amount exceeds balance");
    });

    it("deposits is successful", async function() {
        const before = await token.balanceOf(bridge.address);
        await bridge.depositERC20(
            token.address, 1, ethers.utils.parseEther("1.0")
        );
        const after = await token.balanceOf(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));  
    });

    it("deposits succeed when withdrawals are paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_WITHDRAWALS, nextPauseNonce(), validTimestamp());

        const before = await token.balanceOf(bridge.address);
        await bridge.depositERC20(
            token.address, 1, ethers.utils.parseEther("1.0")
        );
        const after = await token.balanceOf(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));

        await setPaused(bridge, signers, domainSeparator, PAUSE_NOTHING, nextPauseNonce(), validTimestamp());
    });

    it("withdrawals are rejected when bridge is paused", async function() {
        await setPaused(bridge, signers, domainSeparator, PAUSE_WITHDRAWALS, nextPauseNonce(), validTimestamp());

        const recipient = signers[1].address;
        await expect(withdrawERC20(
            bridge,
            token,
            signers,
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

        await expect(
            bridge.depositERC20(token.address, 1, ethers.utils.parseEther("1.0"))
        ).to.be.revertedWith("deposits are paused");

        const recipient = signers[1].address;
        await expect(withdrawERC20(
            bridge,
            token,
            signers,
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
        await expect(withdrawERC20(
            bridge,
            token,
            signers,
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
            token.address,
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawERC201"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawERC20(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with invalid domain separator are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient,
            token.address,
            ethers.utils.parseEther("1.0")
        ];
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version()) + 1, 31337, bridge.address] 
        ));
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
            [invalidDomainSeparator, ethers.utils.id("withdrawERC20"), request]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawERC20(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with invalid signatures are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            token.address,
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawERC20"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        signatures[0] = signatures[1];
        await expect(
            bridge.withdrawERC20(request, signatures, [...Array(20).keys()])
        ).to.be.revertedWith("signature does not match");
    });

    it("withdrawals with non sorted signatures are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient, 
            token.address,
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawERC20"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        const tmp = signatures[1];
        signatures[1] = signatures[0];
        signatures[0] = tmp; 
        const indexes = [...Array(20).keys()];
        indexes[0] = 1;
        indexes[1] = 0;
        await expect(
            bridge.withdrawERC20(request, signatures, indexes)
        ).to.be.revertedWith("signatures not sorted by signer");
    });

    it("withdrawals with indexes length not matching signatures length are rejected", async function() {
        const recipient = signers[1].address;
        const request = [
            ethers.utils.formatBytes32String("0"), 
            validTimestamp(), 
            recipient,
            token.address,
            ethers.utils.parseEther("1.0")
        ];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "tuple(bytes32, uint256, address, address, uint256)"], 
            [domainSeparator, ethers.utils.id("withdrawERC20"), request]
        )));
        const signatures = await Promise.all(signers.slice(0,20).map(s => s.signMessage(hash)));
        await expect(
            bridge.withdrawERC20(request, signatures, [...Array(19).keys()])
        ).to.be.revertedWith("number of signatures does not equal number of indexes");
    });

    it("cannot withdraw more than bridge balance", async function() {
        await expect(
            withdrawERC20(
                bridge,
                token,
                signers,    
                ethers.utils.formatBytes32String("0"),
                domainSeparator,
                validTimestamp(),
                signers[2].address,
                ethers.utils.parseEther("200")
            )
        ).to.be.revertedWith("ERC20: transfer amount exceeds balance");
    });

    it("withdrawal succeeds", async function() {
        let before = await token.balanceOf(bridge.address);
        await bridge.depositERC20(
            token.address, 1, ethers.utils.parseEther("1.0")
        );
        let after = await token.balanceOf(bridge.address);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));  

        const recipient = signers[1].address;
        before = await token.balanceOf(recipient);
        await withdrawERC20(
            bridge,
            token,
            signers,
            ethers.utils.formatBytes32String("0"),
            domainSeparator,
            validTimestamp(),
            recipient,
            ethers.utils.parseEther("1.0")
        );
        after = await token.balanceOf(recipient);
        expect(after.sub(before)).to.equal(ethers.utils.parseEther("1.0"));

        // reusing request id will be rejected
        await expect(
            withdrawERC20(
                bridge,
                token,
                signers,    
                ethers.utils.formatBytes32String("0"),
                domainSeparator,
                validTimestamp(),
                signers[2].address,
                ethers.utils.parseEther("2.0")
            )
        ).to.be.revertedWith("request is already fulfilled");
    });

    it("updateSigners invalidates withdrawal transactions", async function() {
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);
        await expect(
            withdrawERC20(
                bridge,
                token,
                signers,    
                ethers.utils.formatBytes32String("1"),
                domainSeparator,
                validTimestamp(),
                signers[2].address,
                ethers.utils.parseEther("1.0")
            )
        ).to.be.revertedWith("signature does not match");
    });
});

module.exports = { withdrawERC20 };