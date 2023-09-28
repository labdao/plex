import React, { useContext, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setIsLoggedIn, selectIsLoggedIn } from '@/lib/redux';
import { Web3AuthContext } from './Web3AuthContext';  // Adjust the import to your project's structure

const LoginComponent: React.FC = () => {
  const web3Auth = useContext(Web3AuthContext);
  const dispatch = useDispatch();
  const isLoggedIn = useSelector(selectIsLoggedIn)

  useEffect(() => {
    console.log(isLoggedIn);
  }, [isLoggedIn])

  const handleLogin = async () => {
    console.log('handleLogin called');
    if (web3Auth) {
      try {
        console.log(web3Auth)
        const result = await web3Auth.connect();
        console.log(result);
        dispatch(setIsLoggedIn(true));
        console.log(isLoggedIn)
      } catch (error) {
        console.log(error);
      }
    }
  };

  return (
    <div>
      <button onClick={handleLogin}>Login</button>
    </div>
  );
};

export default LoginComponent;