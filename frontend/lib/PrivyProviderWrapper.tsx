'use client';

import { PrivyProvider, User } from '@privy-io/react-auth';
import React, { useState, useEffect, useCallback } from 'react';
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

        try {
            // First, try to fetch the user data
            await dispatch(fetchUserDataAsync()).unwrap();
            console.log('Existing user found, reloading page');
            window.location.reload();
        } catch (error) {
            console.log('User not found, creating new user');
            try {
                // If user not found, create a new user
                await dispatch(saveUserAsync({ walletAddress: user.wallet.address })).unwrap();
                console.log('New user created, fetching user data');
                await dispatch(fetchUserDataAsync()).unwrap();
                console.log('User data fetched successfully, reloading page');
                window.location.reload();
            } catch (saveError) {
                console.error('Error creating or fetching new user:', saveError);
            }
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