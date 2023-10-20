'use client';

import React, { useState } from 'react';
import { PrivyProvider, User } from '@privy-io/react-auth';
import { PrivyAuthContext } from './PrivyContext';
import { optimismGoerli } from '@wagmi/chains';

export default function PrivyProviderWrapper({
    children,   
}: {
    children: React.ReactNode;
}) {
    const [user, setUser] = useState<User | null>(null);
    const [authenticated, setAuthenticated] = useState<boolean>(false);

    const handleLogin = () => {
        setUser(user);
        setAuthenticated(true);
    }

    return (
        <PrivyAuthContext.Provider value={{ user, authenticated }}>
            <PrivyProvider
                appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID || ''}
                onSuccess={handleLogin}
                config={{
                    appearance: {
                        theme: "dark",
                        accentColor: "#6bdaad",
                        logo: "https://imgur.com/6egHxy0.png"
                    },
                    defaultChain: optimismGoerli,
                    supportedChains: [optimismGoerli]
                }}
            >
                {children}
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}