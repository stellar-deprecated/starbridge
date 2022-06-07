const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Bridge", function() {
    it("deploy Bridge contract with three signers", async function() {
      const signers = (await ethers.getSigners()).slice(0, 3);

      const Bridge = await ethers.getContractFactory("Bridge");
      const addresses = signers.map(a => a.address);
      addresses.sort();

      const bridge = await Bridge.deploy(addresses, 2);
      for(let i = 0; i < addresses.length; i++) {
        expect(await bridge.signers(i)).to.equal(addresses[i]);
      }
      await expect(bridge.signers(addresses.length)).to.be.reverted;
    });
  });