'use client';

import { PrivyProvider, User } from '@privy-io/react-auth';
import React, { useState, useCallback } from 'react';
import { optimismSepolia } from 'viem/chains'
import { useDispatch } from 'react-redux';

import { PrivyAuthContext } from './PrivyContext';
import { fetchUserDataAsync, saveUserAsync } from './redux';
import { ReduxDispatch } from '@/lib/redux/store';

export default function PrivyProviderWrapper({
    children,   
}: {
    children: React.ReactNode;
}) {
    const [user, setUser] = useState<User | null>(null);
    const [authenticated, setAuthenticated] = useState<boolean>(false);
    const dispatch = useDispatch<ReduxDispatch>();

    const handleLogin = useCallback(async (user: User) => {
        console.log('Login successful, user:', user);
        setUser(user);
        setAuthenticated(true);
    
        const walletAddress = user.wallet?.address;
        if (!walletAddress) {
            console.error('No wallet address found for user:', user);
            return;
        }
    
        try {
            // This will create the user if they don't exist, or return existing user data
            const userData = await dispatch(saveUserAsync({ walletAddress })).unwrap();
            console.log('User data saved/retrieved:', userData);
            window.location.reload();
        } catch (error) {
            console.error('Error saving/retrieving user data:', error);
        }
    }, [dispatch]);

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