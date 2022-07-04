const { ethers } = require("hardhat");
const { BigNumber } = require("@ethersproject/bignumber");


// run `npx hardhat run scripts/deploy.js --network localhost` to execute this script
async function main() {
    const signers = (await ethers.getSigners()).slice(0, 3);
    signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
    const addresses = signers.map(a => a.address);
    console.log("validators: ", addresses);

    const Bridge = await ethers.getContractFactory("Bridge");

    const bridge = await Bridge.deploy(addresses, 2);
    console.log("Bridge address:", bridge.address);
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
  