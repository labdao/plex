import React, { useEffect, useState, ReactNode } from 'react';
import { Web3Auth } from '@web3auth/modal';
import { Web3AuthContext } from './Web3AuthContext';
import { WALLET_ADAPTERS } from '@web3auth/base';

interface Web3AuthProviderProps {
    children: ReactNode;
}

export const Web3AuthProvider: React.FC<Web3AuthProviderProps> = ({ children }) => {
    const [web3AuthInstance, setWeb3AuthInstance] = useState<Web3Auth | null>(null);
  
    useEffect(() => {
      if (!web3AuthInstance) {
        const instance = new Web3Auth({
            clientId: "BKURHIghKRSWvu0c2IM8hrFtKRny4zLjBqO8Mr4fiIoGc2cSB8_il38d2T5fIxDvpIFQqyEZ5lbNswl_GITUZd0",
            web3AuthNetwork: "sapphire_devnet",
            chainConfig: {
              chainId: "0x1a4",
              chainNamespace: "other",
              rpcTarget: "https://opt-goerli.g.alchemy.com/v2/tvdEoAYqtbNgXzL-Dma7if3i3NzUS3-N",
            },
            uiConfig: {
              loginMethodsOrder: ["email_passwordless"],
              loginGridCol: 2,
              primaryButton: "emailLogin",
            },
        });
        instance.initModal({
            modalConfig: {
                [WALLET_ADAPTERS.OPENLOGIN]: {
                    label: "openlogin",
                    loginMethods: {
                        email_passwordless: {
                            name: "email_passwordless",
                            showOnModal: true,
                        },
                        google: {
                            name: "google",
                            showOnModal: true,
                        },
                        github: {
                            name: "github",
                            showOnModal: true,
                        },
                        facebook: {
                            name: "facebook",
                            showOnModal: false,
                        },
                        discord: {
                            name: "discord",
                            showOnModal: false,
                        },
                        twitch: {
                            name: "twitch",
                            showOnModal: false,
                        },
                        apple: {
                            name: "apple",
                            showOnModal: false,
                        },
                        reddit: {
                            name: "reddit",
                            showOnModal: false,
                        },
                        line: {
                            name: "line",
                            showOnModal: false,
                        },
                        wechat: {
                            name: "wechat",
                            showOnModal: false,
                        },
                        twitter: {
                            name: "twitter",
                            showOnModal: false,
                        },
                        kakao: {
                            name: "kakao",
                            showOnModal: false,
                        },
                        linkedin: {
                            name: "linkedin",
                            showOnModal: false,
                        },
                        weibo: {
                            name: "weibo",
                            showOnModal: false,
                        },
                        sms_passwordless: {
                            name: "sms_passwordless",
                            showOnModal: false,
                        },
                    },
                },
                [WALLET_ADAPTERS.METAMASK]: {
                    label: "metamask",
                    showOnModal: true,
                }
            }
        }).then(() => {
          setWeb3AuthInstance(instance);
        });
      }
    }, []);
  
    return (
      <Web3AuthContext.Provider value={web3AuthInstance}>
        {children}
      </Web3AuthContext.Provider>
    );
  };