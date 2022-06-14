require("@nomiclabs/hardhat-waffle");
require('solidity-coverage');

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: {
    version: "0.8.4",
    settings: {
      optimizer: {
        enabled: true,
        runs: 20000
      },
    }
  },
  networks: {
    goerli: {
      url: "https://ethereum-goerli-rpc.allthatnode.com/",
      accounts: [
        "51138e68e8a5fa906d38c5b42bc01b805d7adb3fce037743fb406bb10aa83307",
        "cff41ce3c1708e589b87198c9ee494eef407ca2a765a4353cf162c85ddc81cd9",
        "0b1037a08795be0955e39e7e279e0530eb89e0ec06d372ff6f122a5a4e1a6f84",
      ]
    }
  }
};
