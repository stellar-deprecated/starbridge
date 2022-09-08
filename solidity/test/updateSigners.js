const { expect } = require("chai");
const { ethers } = require("hardhat");


async function updateSigners(bridge, signers, domainSeparator, newAddresses, newMinThreshold) {
    const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "address[]", "uint8"],
        [domainSeparator, ethers.utils.id("updateSigners"), newAddresses, newMinThreshold]
    )));
    const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
    return bridge.updateSigners(newAddresses, newMinThreshold, signatures, [...Array(signers.length).keys()]);
}

describe("updateSigners", function() {
    let signers;
    let bridge;

    this.beforeAll(async function() {
        signers = (await ethers.getSigners()).slice(0, 20);
        signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = signers.map(a => a.address);

        const Bridge = await ethers.getContractFactory("Bridge");
        bridge = await Bridge.deploy(addresses, 20);
    });

    it("rejects invalid minThreshold values", async function() {
        const addresses = signers.map(a => a.address);
        let domainSeparator = await bridge.domainSeparator();

        for (let i = 0; i <= 10; i++) {
            await expect(updateSigners(bridge, signers, domainSeparator, addresses, i)).to.be.revertedWith("min threshold is too low");
        }
        await expect(updateSigners(bridge, signers, domainSeparator, addresses, 21)).to.be.revertedWith("min threshold is too high");

        let version = await bridge.version();
        for (let i = 11; i <= 20; i++) {
            await updateSigners(bridge, signers, domainSeparator, addresses, i);
            version++;
            expect(await bridge.minThreshold()).to.equal(i);
            expect(await bridge.version()).to.equal(version);
            expect(await bridge.domainSeparator()).to.not.equal(domainSeparator);
            domainSeparator = await bridge.domainSeparator()
        }
    });

    it("rejects invalid method id", async function() {
        const addresses = signers.map(a => a.address);
        const domainSeparator = await bridge.domainSeparator();

        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "address[]", "uint8"],
            [domainSeparator, ethers.utils.id("updateSigners1"), addresses, 20]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(bridge.updateSigners(
            addresses, 20, signatures, [...Array(signers.length).keys()]
        )).to.be.revertedWith("signature does not match");
    });

    it("rejects invalid domain separator", async function() {
        const addresses = signers.map(a => a.address);
        const invalidDomainSeparator = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["uint256", "uint256", "address"],
            [(await bridge.version()) + 1, 31337, bridge.address] 
        ));
        const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "bytes32", "address[]", "uint8"],
            [invalidDomainSeparator, ethers.utils.id("updateSigners"), addresses, 20]
        )));
        const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
        await expect(bridge.updateSigners(
            addresses, 20, signatures, [...Array(signers.length).keys()]
        )).to.be.revertedWith("signature does not match");
    });

    it("rejects invalid signers", async function() {
        const newSigners = [];
        for (let i = 0; i < 256; i++) {
            newSigners.push(ethers.Wallet.createRandom())
        }
        newSigners.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = newSigners.map(a => a.address);
        const domainSeparator = await bridge.domainSeparator();

        await expect(updateSigners(
            bridge, signers, domainSeparator, [addresses[0], addresses[0], addresses[1]], 3
        )).to.be.revertedWith("signers not sorted");
        await expect(updateSigners(
            bridge, signers, domainSeparator, [addresses[0], addresses[1], addresses[2], addresses[4], addresses[3]], 5
        )).to.be.revertedWith("signers not sorted");

        await expect(updateSigners(bridge, signers, domainSeparator, [], 0)).to.be.revertedWith("too few signers");
        await expect(updateSigners(bridge, signers, domainSeparator, addresses, 255)).to.be.revertedWith("too many signers");
    });

    it("succeeds", async function() {
        const newSigners = [];
        for (let i = 0; i < 10; i++) {
            newSigners.push(ethers.Wallet.createRandom())
        }
        newSigners.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
        const addresses = newSigners.map(a => a.address);
        const domainSeparator = await bridge.domainSeparator();

        await updateSigners(bridge, signers, domainSeparator, addresses, addresses.length);

        // replay prevention
        await expect(updateSigners(
            bridge, signers, domainSeparator, addresses, addresses.length
        )).to.be.revertedWith("signature does not match");

        for(let i = 0; i < addresses.length; i++) {
            expect(await bridge.signers(i)).to.equal(addresses[i]);
        }
        await expect(bridge.signers(addresses.length)).to.be.reverted;
    });
});

module.exports = { updateSigners };
