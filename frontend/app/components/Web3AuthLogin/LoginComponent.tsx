import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setIsLoggedIn, selectIsLoggedIn } from '@/lib/redux';
import { Web3Auth } from '@web3auth/modal';
import { AppDispatch, ReduxState } from '@/lib/redux/store'; // Import the RootState from your store

import Button from '@mui/material/Button';
import Box from '@mui/material/Box';

const LoginComponent: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const [web3AuthInstance, setWeb3AuthInstance] = useState<Web3Auth | null>(null);

  useEffect(() => {
    if (!web3AuthInstance) {
      const instance = new Web3Auth({
        clientId: "BKURHIghKRSWvu0c2IM8hrFtKRny4zLjBqO8Mr4fiIoGc2cSB8_il38d2T5fIxDvpIFQqyEZ5lbNswl_GITUZd0",
        web3AuthNetwork: "sapphire_devnet",
        chainConfig: {
          chainId: "0x1a4",
          chainNamespace: "other",
          rpcTarget: "https://opt-goerli.g.alchemy.com/v2/tvdEoAYqtbNgXzL-Dma7if3i3NzUS3-N",
        },
        uiConfig: {
          loginMethodsOrder: ["email_passwordless"],
        },
      });
      instance.initModal().then(() => {
        console.log("modal initialized")
        setWeb3AuthInstance(instance);
      });
    }
  }, []);

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
      } catch (error) {
        console.error(error);
      }
    }
  };

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