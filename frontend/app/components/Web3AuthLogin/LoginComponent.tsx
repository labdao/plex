import React, { useContext, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setIsLoggedIn, selectIsLoggedIn, setUsername, setWalletAddress } from '@/lib/redux';
import { AppDispatch, ReduxState } from '@/lib/redux/store'; // Import the RootState from your store

import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import { Web3AuthContext } from '../../../lib/Web3AuthContext'
import { publicToAddress } from 'ethereumjs-util';

const LoginComponent: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const web3AuthInstance = useContext(Web3AuthContext);

  const handleLogin = async () => {
    console.log("handling login")
    if (web3AuthInstance) {
      try {
        console.log(web3AuthInstance)
        const result = await web3AuthInstance.connect();
        console.log(result)
        if (result) {
          dispatch(setIsLoggedIn(true));
        }
        getUserInfo();
        getWalletAddress();
      } catch (error) {
        console.error(error);
      }
    }
  };

  const getUserInfo = async () => {
    try {
      if (web3AuthInstance) {
        const response = await web3AuthInstance.getUserInfo();
        console.log(response);
        const email = response.email as string;
        dispatch(setUsername(email))
      }
    } catch (error) {
      console.error("Failed to get user info:", error);
    }
  }

  const getWalletAddress = async () => {
    try {
      if (web3AuthInstance) {
        const response = await web3AuthInstance.authenticateUser() as any;
        console.log(response);
        const wallet = response.wallets[0];
        const addressBuffer = publicToAddress(Buffer.from(wallet.public_key, "hex"), true);
        const address = `0x${addressBuffer.toString("hex")}`;
        console.log("Wallet address:", address);
        dispatch(setWalletAddress(address));
      }
    } catch (error) {
      console.error("Failed to get wallet address:", error);
    }
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
  );
};

export default LoginComponent;