<<<<<<< HEAD
import { useLogin, usePrivy } from "@privy-io/react-auth";
=======
import { useLogin } from "@privy-io/react-auth";
import { usePrivy } from "@privy-io/react-auth";
import { Slot } from "@radix-ui/react-slot";
>>>>>>> main
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
      if (!walletAddress) {
        console.log("No wallet address found");
        return;
      }

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
