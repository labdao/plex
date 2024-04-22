// import { switchChain } from "@privy-io/react-auth";
import { useLogin, useWallets } from "@privy-io/react-auth";
import { usePrivy } from "@privy-io/react-auth";
import { Slot } from "@radix-ui/react-slot";
import React from "react";
import { useDispatch } from "react-redux";

import { ButtonProps } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import { AppDispatch } from "@/lib/redux";
import { saveUserAsync } from "@/lib/redux/slices/userSlice/thunks";

const PrivyLoginButton = (props: ButtonProps) => {
  const dispatch: AppDispatch = useDispatch();
  const { ready, authenticated } = usePrivy();
  const Comp = props.asChild ? Slot : Button;

  //This component must remain mounted wherever you use it for the callback to fire correctly
  const { login } = useLogin({
    onComplete: (user, isNewUser, wasAlreadyAuthenticated) => {
      const walletAddress = user?.wallet?.address;
      const chainId = user?.wallet?.chainId;

      if (!walletAddress) {
        console.log("No wallet address found");
        return;
      }

      // if (chainId !== 'eip155:11155420') {
      //   console.log(`User is on a different network: ${chainId}. Switching to OP Sepolia testnet.`);
      //   try {
      //     await switchChain('eip155:11155420');
      //     console.log("Switched to OP Sepolia testnet successfully.");
      //   } catch (switchError) {
      //     console.error("Failed to switch to OP Sepolia testnet:", switchError);
      //   }
      // } else {
      //   console.log("User is already on OP Sepolia testnet.");
      // }

      console.log(`User authenticated with wallet address: ${walletAddress} on chainId: ${chainId}`)

      if (isNewUser) {
        console.log("New user");
      } else if (wasAlreadyAuthenticated) {
        console.log("User was already authenticated");
      } else {
        console.log("User authenticated");
      }
      dispatch(saveUserAsync({ walletAddress }));
    },
    onError: (error) => {
      console.log("onError callback triggered", error);
    },
  });

  const handleLogin = async () => {
    try {
      login();
    } catch (error) {
      console.log("Error calling login function:", error);
    }
  };

  return ready && !authenticated ? <Comp onClick={handleLogin} {...props} /> : null;
};

export default PrivyLoginButton;
