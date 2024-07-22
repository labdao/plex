'use client';

import { PrivyProvider, User } from '@privy-io/react-auth';
import React, { useState, useEffect } from 'react';
import { optimismSepolia } from 'viem/chains'
import { useDispatch } from 'react-redux';

import { PrivyAuthContext } from './PrivyContext';
import { fetchUserDataAsync } from './redux';
import { ReduxDispatch } from '@/lib/redux/store';

export default function PrivyProviderWrapper({
    children,   
}: {
    children: React.ReactNode;
}) {
    const [user, setUser] = useState<User | null>(null);
    const [authenticated, setAuthenticated] = useState<boolean>(false);
    const dispatch = useDispatch<ReduxDispatch>();

    useEffect(() => {
        // Check if this is a fresh login
        const isFirstLogin = sessionStorage.getItem('isFirstLogin') === 'true';
        if (isFirstLogin) {
            sessionStorage.removeItem('isFirstLogin');
            dispatch(fetchUserDataAsync())
                .unwrap()
                .then(() => {
                    // Use window.location.reload() instead of router.reload()
                    window.location.reload();
                })
                .catch((error: string) => {
                    console.error('Error fetching user data:', error);
                    // Handle error (e.g., show an error message to the user)
                });
        }
    }, [dispatch]);

    const handleLogin = (user: User) => {
        setUser(user);
        setAuthenticated(true);
        // Set the first login flag
        sessionStorage.setItem('isFirstLogin', 'true');
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
                    defaultChain: optimismSepolia,
                    supportedChains: [optimismSepolia]
                }}
            >
                {children}
            </PrivyProvider>
        </PrivyAuthContext.Provider>
    )
}