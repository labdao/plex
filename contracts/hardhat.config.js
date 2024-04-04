require("@nomicfoundation/hardhat-toolbox");
require('dotenv').config();

/** @type import('hardhat/config').HardhatUserConfig */

const ALCHEMY_API_KEY = process.env.ALCHEMY_API_KEY;
const WALLET_PRIVATE_KEY = process.env.WALLET_PRIVATE_KEY;
const OPTIMISM_BLOCK_EXPLORER_API_KEY = process.env.OPTIMISM_BLOCK_EXPLORER_API_KEY;

module.exports = {
  solidity: "0.8.20",
  networks: {
    'optimism-sepolia': {
      url: `https://opt-sepolia.g.alchemy.com/v2/${ALCHEMY_API_KEY}`,
      accounts: [WALLET_PRIVATE_KEY],
    }
  },
  etherscan: {
    apiKey: OPTIMISM_BLOCK_EXPLORER_API_KEY,
    customChains: [
      {
        network: 'optimism-sepolia',
        chainId: 11155420,
        urls: {
          apiURL: 'https://api-sepolia-optimistic.etherscan.io/api',
        }
      }
    ]
  }
};
