
const { expect } = require("chai");
const { ethers } = require("hardhat");
const { updateSigners } = require("./updateSigners");
const { withdrawERC20 } = require("./erc20");
const { validTimestamp } = require("./util");


describe("registerStellarAsset", function() {
    let signers;
    let bridge;
    let wrappedXLM;
    let recipient;

    async function registerStellarAsset(domainSeparator, name, symbol, decimals) {
        const request = [decimals, name, symbol];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                domainSeparator, 
                ethers.utils.id("registerStellarAsset"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        return await bridge.registerStellarAsset(request, signatures, [...Array(20).keys()]);
    }

    async function getToken(tx) {
        const result = await tx.wait();
        const events = await bridge.queryFilter(bridge.filters.RegisterStellarAsset(), result.blockNumber, result.blockNumber);
        expect(events.length).to.be.eql(1);
        return ethers.getContractAt("StellarAsset", events[0].args[0]);
    }

    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        recipient = signers[0];
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);

        wrappedXLM = await getToken(await registerStellarAsset(await bridge.domainSeparator(), "Stellar Lumens", "XLM", 7));
        expect(await wrappedXLM.decimals()).to.be.eql(7);
        expect(await wrappedXLM.name()).to.be.eql("Stellar Lumens");
        expect(await wrappedXLM.symbol()).to.be.eql("XLM");
        expect(await bridge.isStellarAsset(wrappedXLM.address)).to.be.true;
    });

    it("rejects duplicate transactions", async function() {
        let domainSeparator = await bridge.domainSeparator();
        await expect(registerStellarAsset(domainSeparator, "Stellar Lumens", "XLM", 7)).to.be.reverted;
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);
        domainSeparator = await bridge.domainSeparator();
        await expect(registerStellarAsset(domainSeparator, "Stellar Lumens", "XLM", 7)).to.be.reverted;
    });

    it("updateSigners invalidates transactions", async function() {
        let domainSeparator = await bridge.domainSeparator();
        await updateSigners(bridge, signers, domainSeparator, signers.map(s => s.address), signers.length);
        await expect(registerStellarAsset(domainSeparator, "wrapped yXLM", "yXLM", 7)).to.be.reverted;
        domainSeparator = await bridge.domainSeparator();
        await expect(registerStellarAsset(domainSeparator, "wrapped yXLM", "yXLM", 7)).to.not.be.reverted;
    });

    it("StellarAssets.mint() cannot be called", async function() {
        await expect(wrappedXLM.mint(signers[0].address, ethers.utils.parseEther("100.0"))).to.be.reverted;
    });

    it("withdrawERC20 and depositERC20 mints and burns StellarAsset tokens", async function() {
        expect(await wrappedXLM.balanceOf(recipient.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.balanceOf(bridge.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.totalSupply()).to.be.equal(ethers.utils.parseEther("0"));

        await withdrawERC20(
            bridge,
            wrappedXLM,
            signers,
            ethers.utils.formatBytes32String("0"),
            await bridge.domainSeparator(),
            validTimestamp(),
            recipient.address,
            ethers.utils.parseEther("3.0"),
        );

        expect(await wrappedXLM.balanceOf(recipient.address)).to.be.equal(ethers.utils.parseEther("3.0"));
        expect(await wrappedXLM.balanceOf(bridge.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.totalSupply()).to.be.equal(ethers.utils.parseEther("3.0"));

        await bridge.depositERC20(
            wrappedXLM.address, 1, ethers.utils.parseEther("1.0")
        );

        expect(await wrappedXLM.balanceOf(recipient.address)).to.be.equal(ethers.utils.parseEther("2.0"));
        expect(await wrappedXLM.balanceOf(bridge.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.totalSupply()).to.be.equal(ethers.utils.parseEther("2.0"));

        await bridge.depositERC20(
            wrappedXLM.address, 1, ethers.utils.parseEther("2.0")
        );

        expect(await wrappedXLM.balanceOf(recipient.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.balanceOf(bridge.address)).to.be.equal(ethers.utils.parseEther("0"));
        expect(await wrappedXLM.totalSupply()).to.be.equal(ethers.utils.parseEther("0"));
    });

    it("cannot deposit more tokens than current balance", async function() {
        expect(await wrappedXLM.balanceOf(recipient.address)).to.be.equal(ethers.utils.parseEther("0"));
        await expect(bridge.depositERC20(
            wrappedXLM.address, 1, ethers.utils.parseEther("2.0")
        )).revertedWith("ERC20: burn amount exceeds balance");
    });

    it("transactions with invalid method id are rejected", async function() {
        const decimals = 7;
        const name = "Wrapped USDC";
        const symbol = "USDC";
        const request = [decimals, name, symbol];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                await bridge.domainSeparator(), 
                ethers.utils.id("registerStellarAsset1"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.registerStellarAsset(request, signatures, [...Array(20).keys()])
        ).revertedWith("signature does not match");
    });

    it("transactions with invalid domain separator are rejected", async function() {
        const decimals = 7;
        const name = "Wrapped USDC";
        const symbol = "USDC";
        const request = [decimals, name, symbol];
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version()) + 1, 31337, bridge.address] 
        ));
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                invalidDomainSeparator, 
                ethers.utils.id("registerStellarAsset"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.registerStellarAsset(request, signatures, [...Array(20).keys()])
        ).revertedWith("signature does not match");
    });

    it("transactions with invalid signatures are rejected", async function() {
        const decimals = 7;
        const name = "Wrapped USDC";
        const symbol = "USDC";
        const request = [decimals, name, symbol];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                await bridge.domainSeparator(), 
                ethers.utils.id("registerStellarAsset"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        signatures[0] = signatures[1];
        await expect(
            bridge.registerStellarAsset(request, signatures, [...Array(20).keys()])
        ).revertedWith("signature does not match");
    });

    it("transactions with non sorted signatures are rejected", async function() {
        const decimals = 7;
        const name = "Wrapped USDC";
        const symbol = "USDC";
        const request = [decimals, name, symbol];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                await bridge.domainSeparator(), 
                ethers.utils.id("registerStellarAsset"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        const tmp = signatures[1];
        signatures[1] = signatures[0];
        signatures[0] = tmp; 
        const indexes = [...Array(20).keys()];
        indexes[0] = 1;
        indexes[1] = 0;
        await expect(
            bridge.registerStellarAsset(request, signatures, indexes)
        ).revertedWith("signatures not sorted by signer");
    });

    it("transactions with indexes length not matching signatures length are rejected", async function() {
        const decimals = 7;
        const name = "Wrapped USDC";
        const symbol = "USDC";
        const request = [decimals, name, symbol];
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "uint8", "bytes32", "bytes32"], 
            [
                await bridge.domainSeparator(), 
                ethers.utils.id("registerStellarAsset"), 
                decimals, 
                ethers.utils.id(name),
                ethers.utils.id(symbol),
            ]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(
            bridge.registerStellarAsset(request, signatures, [...Array(19).keys()])
        ).to.be.revertedWith("number of signatures does not equal number of indexes");
    });
});