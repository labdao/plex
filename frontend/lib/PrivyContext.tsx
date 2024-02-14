import { User } from '@privy-io/react-auth';
import React from 'react';

interface AuthState {
    user: User | null;
    authenticated: boolean;
}

export const PrivyAuthContext = React.createContext<AuthState>({ user: null, authenticated: false});