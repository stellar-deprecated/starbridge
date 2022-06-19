const { ethers } = require("hardhat");


// run `npx hardhat run scripts/deploy.js --network goerli` to execute this script
async function main() {
    const signers = (await ethers.getSigners());
    signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
    const addresses = signers.map(a => a.address);

    const Bridge = await ethers.getContractFactory("Bridge");

    console.log("validators: ", addresses);
    const bridge = await Bridge.deploy(addresses, addresses.length);
  
    console.log("Bridge address:", bridge.address);
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
  