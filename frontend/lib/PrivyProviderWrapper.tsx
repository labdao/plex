'use client';

import React, { useState } from 'react';
import { PrivyProvider, User } from '@privy-io/react-auth';
import { PrivyAuthContext } from './PrivyContext';
import { PrivyWagmiConnector } from '@privy-io/wagmi-connector';
import { alchemyProvider } from 'wagmi/providers/alchemy';
import { optimismGoerli } from '@wagmi/chains';
import { configureChains } from 'wagmi';

export default function PrivyProviderWrapper({
    children,   
}: {
    children: React.ReactNode;
}) {
    const [user, setUser] = useState<User | null>(null);
    const [authenticated, setAuthenticated] = useState<boolean>(false);

    const ALCHEMY_RPC: string = process.env.NEXT_PUBLIC_OPTIMISM_GOERLI_ALCHEMY_RPC || '';
    const configureChainsConfig = configureChains(
        [optimismGoerli],
        [alchemyProvider({ apiKey: ALCHEMY_RPC})]
    )

    const handleLogin = () => {
        setUser(user);
        setAuthenticated(true);
    }

    return (
        <PrivyAuthContext.Provider value={{ user, authenticated }}>
            <PrivyProvider
                appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID || 'clo7adk6w07q7jq0f08yrnkur'}
                onSuccess={handleLogin}
                config={{
                    appearance: {
                        theme: "dark",
                        accentColor: "#6bdaad",
                        logo: "https://imgur.com/6egHxy0.png"
                    },
                    // defaultChain: optimismGoerli,
                    // supportedChains: [optimismGoerli]
                }}
            >
                <PrivyWagmiConnector wagmiChainsConfig={configureChainsConfig}>
                    {children}
                </PrivyWagmiConnector>
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}