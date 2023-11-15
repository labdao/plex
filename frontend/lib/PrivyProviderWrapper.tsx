'use client';

import { PrivyProvider, User } from '@privy-io/react-auth';
import React, { useState } from 'react';

import { PrivyAuthContext } from './PrivyContext';

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
                appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID || 'clo7adk6w07q7jq0f08yrnkur'}
                onSuccess={handleLogin}
                config={{
                    appearance: {
                        theme: "dark",
                        accentColor: "#6bdaad",
                    }
                }}
            >
                {children}
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}