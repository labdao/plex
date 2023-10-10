'use client';

import React, { useState } from 'react';
import { PrivyProvider, User } from '@privy-io/react-auth';
import { PrivyAuthContext } from './PrivyContext';
import { useRouter } from 'next/router';

export default function PrivyProviderWrapper({
    children,   
}: {
    children: React.ReactNode;
}) {
    const [user, setUser] = useState<User | null>(null);

    const handleLogin = () => {
        setUser(user);
    }

    return (
        <PrivyAuthContext.Provider value={user}>
            <PrivyProvider
                appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID || ''}
                onSuccess={handleLogin}
                config={{
                    appearance: {
                        theme: "dark",
                        accentColor: "#6bdaad",
                        logo: "https://imgur.com/6egHxy0.png"
                    }
                }}
            >
                {children}
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}