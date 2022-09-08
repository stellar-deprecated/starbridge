const { ethers } = require("hardhat");

async function registerStellarAsset(bridge, signers, configVersion, name, symbol, decimals) {
  const request = [decimals, name, symbol];
  const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
      ["uint256", "bytes32", "uint8", "bytes32", "bytes32"], 
      [
          configVersion, 
          ethers.utils.id("registerStellarAsset"), 
          decimals, 
          ethers.utils.id(name),
          ethers.utils.id(symbol),
      ]
  )));
  const signatures = await Promise.all(signers.map(s => s.signMessage(hash)));
  return await bridge.registerStellarAsset(request, signatures, [...Array(signers.length).keys()]);
}

async function getToken(bridge, tx) {
  const result = await tx.wait();
  const events = await bridge.queryFilter(bridge.filters.RegisterStellarAsset(), result.blockNumber, result.blockNumber);
  return ethers.getContractAt("StellarAsset", events[0].args[0]);
}

// run `npx hardhat run scripts/deploy.js --network localhost` to execute this script
async function main() {
    const signers = (await ethers.getSigners()).slice(0, 3);
    signers.sort((a, b) => a.address.toLowerCase().localeCompare(b.address.toLowerCase()));
    const addresses = signers.map(a => a.address);
    console.log("validators: ", addresses);

    const Bridge = await ethers.getContractFactory("Bridge");

    const bridge = await Bridge.deploy(addresses, 2);
    console.log("Bridge address:", bridge.address);

    const wrappedXLM = await getToken(bridge, await registerStellarAsset(bridge, signers, 0, "Stellar Lumens", "XLM", 7));
    console.log("Wrapped XLM address:", wrappedXLM.address);

    console.log("domain separator: ", await bridge.domainSeparator());
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
  