import React, { createContext, useState, useEffect, ReactNode } from 'react';
import { Web3Auth } from '@web3auth/modal';

interface Web3AuthContextProps {
  children: ReactNode;
}

export const Web3AuthContext = createContext<Web3Auth | null>(null);

export const Web3AuthProvider: React.FC<Web3AuthContextProps> = ({ children }) => {
  const [web3Auth, setWeb3Auth] = useState<Web3Auth | null>(null);

  useEffect(() => {
    const initWeb3Auth = async () => {
      const web3AuthInstance = new Web3Auth({
        clientId: "BKURHIghKRSWvu0c2IM8hrFtKRny4zLjBqO8Mr4fiIoGc2cSB8_il38d2T5fIxDvpIFQqyEZ5lbNswl_GITUZd0",
        web3AuthNetwork: "sapphire_devnet",
        chainConfig: {
          chainId: "0x1a4",
          chainNamespace: "other",
          rpcTarget: "https://opt-goerli.g.alchemy.com/v2/tvdEoAYqtbNgXzL-Dma7if3i3NzUS3-N",
        },
        uiConfig: {
          loginMethodsOrder: ["email_passwordless"],
        },
      });
      await web3AuthInstance.initModal();
      setWeb3Auth(web3AuthInstance);
    };
    initWeb3Auth();
  }, []);

  return (
    <Web3AuthContext.Provider value={web3Auth}>
      {children}
    </Web3AuthContext.Provider>
  );
};