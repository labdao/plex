import React, { useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setIsLoggedIn, selectIsLoggedIn } from '@/lib/redux';
import { AppDispatch, ReduxState } from '@/lib/redux/store'; // Import the RootState from your store

import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import { Web3AuthContext } from '../../../lib/Web3AuthContext'

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