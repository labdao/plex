import React, { useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setIsLoggedIn, selectIsLoggedIn, setWalletAddress } from '@/lib/redux';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import { usePrivy, useWallets } from '@privy-io/react-auth';
import { PrivyAuthContext } from '../../../lib/PrivyContext';

import { useRouter } from 'next/navigation'

const PrivyLoginComponent: React.FC = () => {
    const dispatch = useDispatch();
    const isLoggedIn = useSelector(selectIsLoggedIn);
    const user = useContext(PrivyAuthContext);
    const { wallets } = useWallets();

    const router = useRouter()

    const { login } = usePrivy();
    const { authenticated } = usePrivy();

    const handleLogin = async () => {
        if (!user) {
            try {
                await login();
                const walletAddress = await getWalletAddress();
                if (walletAddress) {
                    // dispatch(saveUserAsync(walletAddress));
                }
            } catch (error) {
                console.log(error);
            }
        }
    };

    // getting embedded wallet address
    const getWalletAddress = async () => {
        // may need to be updated based on how we manage users adding multiple wallets
        const walletAddress = wallets[0].address;
        dispatch(setWalletAddress(walletAddress));
        return walletAddress;
    }

    return (
        <Box
            display="flex"
            justifyContent="center"
            mt={2}
        >
            <Button
                variant="contained"
                onClick={handleLogin}
                sx={{ backgroundColor: '#333333', '&:hover': { backgroundColor: '#6bdaad' } }}
            >
                Login
            </Button>
        </Box>
    )
}

export default PrivyLoginComponent;