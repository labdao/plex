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
                appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID || 'clnf5ptkk01h1jn0fzhh4xldt'}
                onSuccess={handleLogin}
                config={{
                    appearance: {
                        theme: "dark",
                        accentColor: "#6bdaad",
                        logo: "https://raw.githubusercontent.com/labdao/plex/main/LabBio_whitetext_transparent.png",
                    },
                    legal: {
                        termsAndConditionsUrl: "https://docs.google.com/document/d/12-GU8FcXEIvnWlPO6WVjgE0Xyo6kXacL/edit?usp=sharing&ouid=114685327806757138792&rtpof=true&sd=true",
                    }
                }}
            >
                {children}
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}