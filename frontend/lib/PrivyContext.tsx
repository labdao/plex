import React from 'react';
import { User } from '@privy-io/react-auth';

interface AuthState {
    user: User | null;
    authenticated: boolean;
}

// export const PrivyAuthContext = React.createContext<User | null>(null);
export const PrivyAuthContext = React.createContext<AuthState>({ user: null, authenticated: false});