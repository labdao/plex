import { hexValue } from '@ethersproject/bytes';
import { Contract } from '@ethersproject/contracts';
import { JsonRpcProvider, Web3Provider } from '@ethersproject/providers';

const CONTRACT_ADDRESS = '0x7336371ce024De5BA5fd80f53594fe518fb793AE';
const ContractArtifact = require('../../contracts/artifacts/contracts/ProofOfScience.sol/ProofOfScience.json');
const abi = ContractArtifact.abi;

export const getEthereumProvider = () => {
  if (typeof window !== 'undefined' && typeof window.ethereum !== 'undefined') {
    const provider = new Web3Provider(window.ethereum);
    provider.send("wallet_switchEthereumChain", [{ chainId: hexValue(11155420) }])
      .catch((error: any) => {
        if (error.code === 4902) {
          console.error("The Optimism Sepolia network is not available in your MetaMask, please add it manually.");
        }
      });
    return provider;
  } else {
    return new JsonRpcProvider('https://sepolia.optimism.io', {
      chainId: 11155420,
      name: 'OP Sepolia'
    });
  }
};

export const getNFTContract = () => {
  const provider = getEthereumProvider();
  return new Contract(CONTRACT_ADDRESS, abi, provider);
};