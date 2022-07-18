require("@nomiclabs/hardhat-waffle");
require("@nomiclabs/hardhat-etherscan");
require('solidity-coverage');

const testAccounts = [
  { privateKey: "51138e68e8a5fa906d38c5b42bc01b805d7adb3fce037743fb406bb10aa83307", balance: "10000000000000000000000" },
  { privateKey: "cff41ce3c1708e589b87198c9ee494eef407ca2a765a4353cf162c85ddc81cd9", balance: "10000000000000000000000" },
  { privateKey: "0b1037a08795be0955e39e7e279e0530eb89e0ec06d372ff6f122a5a4e1a6f84", balance: "10000000000000000000000" },

  { privateKey: "2c6c7345077d7b96679d5e2cf104ed29be3475860dbd92ac472bf8a61ad6464a", balance: "10000000000000000000000" },
  { privateKey: "69a2b099a902fd5008507d79d93898dca86ba4766d3eaa51c10239b0ac4c0d16", balance: "10000000000000000000000" },
  { privateKey: "4bf22efd1efdeaeea709e0cf641f141dff598fe549b24f349a4629304e20b9fa", balance: "10000000000000000000000" },

  { privateKey: "9781c27721e1423190f289053001670a1c60961d32dcbacd4c9ee00b2df9a88f", balance: "10000000000000000000000" },
  { privateKey: "68d318aa94dcdbc266985a74ed98e38f8b287ac543adffcac57356e33f1c113b", balance: "10000000000000000000000" },
  { privateKey: "d388e47585f151455fed2dd5f9fab2d990d140923d29cf89387657a39736d620", balance: "10000000000000000000000" },

  { privateKey: "c274142ea307b7043e56e1e54d3befb2ab1056e80633e51a8a0958f69ddc0a4d", balance: "10000000000000000000000" },
  { privateKey: "636b7f2c628d2a6b9e742fc590f2ccab1f5f1d7d2061d9ad46353605c74769b9", balance: "10000000000000000000000" },
  { privateKey: "3ef813628f3638a58e77f44c9e8830760fa948e9952323f9b47a0083ee30d226", balance: "10000000000000000000000" },

  { privateKey: "534bac3de21dec2149a2f7529edeac119f17864e83ed8d2a88b249e747c5b36a", balance: "10000000000000000000000" },
  { privateKey: "609c03d127f1e1fed73541c09004adb74b78c2b8dd81502df1f921ae1dd100bd", balance: "10000000000000000000000" },
  { privateKey: "de422bb8c12652446edf3ae049d6f6238d9bc53b33257ba36d3d3173657cce92", balance: "10000000000000000000000" },

  { privateKey: "a602617bd95972ba004e7f02c971d1eb9558d5a6a5a0ce3a189037fd5a7ed63b", balance: "10000000000000000000000" },
  { privateKey: "d0dd4ac71b87c1ee18fdf72b05da607a1b9c94d2e1e5f825a29c3bd48760be0b", balance: "10000000000000000000000" },
  { privateKey: "22352a9be4915aab03e69061e706d76c42571257c70d8d996f35fd35a099e256", balance: "10000000000000000000000" },

  { privateKey: "c84cd9162ea55da7da04d3aa7294b8fcd55d31ba583f7ee63addd8cb4ded8cfe", balance: "10000000000000000000000" },
  { privateKey: "4468af6d18b0fdb71774b7f183706e2743839f22e949846e2d186bd1d0fbf48f", balance: "10000000000000000000000" },
  { privateKey: "be40b9c3815b568777283ad72ddf47a47d7ea4358c5d79dcfc6038d264c55050", balance: "10000000000000000000000" },

  { privateKey: "10d8afcd7ef4a9928b63404805f2d1cebc9ae0645e727217fc8668db690f9249", balance: "10000000000000000000000" },
  { privateKey: "ba75e69b948a9e41c8e1c76e5153657e4c9e74add43418b3988fa69089a7ea28", balance: "10000000000000000000000" },
  { privateKey: "921a893c8228c6b23dc3f116387049441b215567b5cee40cf00a737870b72eb8", balance: "10000000000000000000000" },

  { privateKey: "59bb177fb32d92466eb90c4dc5350788f8d251fed4ccfc915cde133e77688d3a", balance: "10000000000000000000000" },
  { privateKey: "819cb9b6082e22cfca504d1876c6b9add3d8b572ef774a018fecf57c3a103c68", balance: "10000000000000000000000" },
  { privateKey: "103999000832f96f28dcac8781501bea6e426d2ae12e4d1e6c6914840b19b38d", balance: "10000000000000000000000" },

  { privateKey: "201a7b1f0d0f1f81908e23fc38ca531a17b78f1ee27f5ab3beb62be569c8a068", balance: "10000000000000000000000" },
  { privateKey: "c1a4af60400ffd1473ada8425cff9f91b533194d6dd30424a17f356e418ac35b", balance: "10000000000000000000000" },
];

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
    hardhat: {
      accounts: testAccounts,
    },
    docker: {
      url: "http://host.docker.internal:8545",
      accounts: testAccounts.map( e => e.privateKey),
    },
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
