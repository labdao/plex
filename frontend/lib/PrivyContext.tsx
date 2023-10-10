import React from 'react';
import { User } from '@privy-io/react-auth';

export const PrivyAuthContext = React.createContext<User | null>(null);