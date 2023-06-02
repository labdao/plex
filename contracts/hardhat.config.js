require("@nomicfoundation/hardhat-toolbox");
require('dotenv').config();

/** @type import('hardhat/config').HardhatUserConfig */

const ALCHEMY_API_KEY = process.env.ALCHEMY_API_KEY;
const WALLET_PRIVATE_KEY = process.env.WALLET_PRIVATE_KEY;
const OPTIMISM_BLOCK_EXPLORER_API_KEY = process.env.OPTIMISM_BLOCK_EXPLORER_API_KEY;

module.exports = {
  solidity: "0.8.18",
  networks: {
    optimismGoerli: {
      url: `https://opt-goerli.g.alchemy.com/v2/${ALCHEMY_API_KEY}`,
      accounts: [WALLET_PRIVATE_KEY]
    }
  },
  etherscan: {
    apiKey: OPTIMISM_BLOCK_EXPLORER_API_KEY
  }
};
