const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Bridge", function() {
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

      const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(["uint256", "bytes32", "address[]", "uint8"], [0, ethers.utils.id("updateSigners"), addresses, 20])));
      const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
      const gas = await bridge.estimateGas.updateSigners(addresses, 20, signatures, [...Array(20).keys()]);
      console.log(gas);
    });
  });