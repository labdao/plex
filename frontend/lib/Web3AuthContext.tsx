import React from 'react';
import { Web3Auth } from '@web3auth/modal';

export const Web3AuthContext = React.createContext<Web3Auth | null>(null);