const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Deploy Bridge", function() {
    it("deploy Bridge contract with 20 signers", async function() {
      const signers = (await ethers.getSigners()).slice(0, 20);
      signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));

      const Bridge = await ethers.getContractFactory("Bridge");
      const addresses = signers.map(a => a.address);

      const bridge = await Bridge.deploy(addresses, 20);
      for(let i = 0; i < 20; i++) {
        expect(await bridge.signers(i)).to.equal(addresses[i]);
      }
      await expect(bridge.signers(20)).to.be.reverted;
    });

    it("deploy Bridge contract with invalid minThreshold", async function() {
      const signers = (await ethers.getSigners()).slice(0, 20);
      signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
      const Bridge = await ethers.getContractFactory("Bridge");
      const addresses = signers.map(a => a.address);

      for (let i = 0; i <= 10; i++) {
        await expect(Bridge.deploy(addresses, i)).to.be.revertedWith("min threshold is too low");
      }

      let domainSeparator = '';
      for (let i = 11; i <= 20; i++) {
        const bridge = await Bridge.deploy(addresses, i);
        expect(await bridge.minThreshold()).to.equal(i);
        expect(await bridge.version()).to.equal(0);
        expect(await bridge.domainSeparator()).to.not.equal(domainSeparator);
        domainSeparator = await bridge.domainSeparator()
      }

      await expect(Bridge.deploy(addresses, 21)).to.be.revertedWith("min threshold is too high");
    });

    it("deploy Bridge contract with invalid signers", async function() {
      const signers = [];
      for (let i = 0; i < 256; i++) {
        signers.push(ethers.Wallet.createRandom())
      }
      signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));


      const Bridge = await ethers.getContractFactory("Bridge");
      const addresses = signers.map(a => a.address);

      await expect(Bridge.deploy([], 0)).to.be.revertedWith("too few signers");
      await expect(Bridge.deploy(addresses, 255)).to.be.revertedWith("too many signers");
      await expect(Bridge.deploy(addresses.slice(0, 255), 255)).to.not.be.reverted;
      await expect(Bridge.deploy([addresses[0], addresses[0], addresses[1]], 3)).to.be.revertedWith("signers not sorted");
      await expect(Bridge.deploy([addresses[0], addresses[1], addresses[2], addresses[4], addresses[3]], 5)).to.be.revertedWith("signers not sorted");
    });
});