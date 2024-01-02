import { useLogin, useWallets } from "@privy-io/react-auth";
import { PrivyAuthContext } from "lib/PrivyContext";
import { useRouter } from "next/navigation";
import React, { useContext } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Button } from "@/components/ui/button";
import { AppDispatch, selectIsLoggedIn, setIsLoggedIn, setWalletAddress } from "@/lib/redux";
import { saveUserAsync } from "@/lib/redux/slices/userSlice/thunks";

const PrivyLogin: React.FC = () => {
  const dispatch: AppDispatch = useDispatch();
  const { user } = useContext(PrivyAuthContext);
  const router = useRouter();

  const { login } = useLogin({
    onComplete: async (user, isNewUser, wasAlreadyAuthenticated) => {
      const walletAddress = await getWalletAddress();
      if (wasAlreadyAuthenticated) {
        console.log("User was already authenticated");
        dispatch(setIsLoggedIn(true));
        router.push("/");
      } else if (isNewUser) {
        console.log("New user");
        dispatch(saveUserAsync({ walletAddress }));
        dispatch(setIsLoggedIn(true));
        router.push("/");
      } else if (user) {
        console.log("User authenticated");
        dispatch(setIsLoggedIn(true));
        router.push("/");
      }
    },
    onError: (error) => {
      console.log("onError callback triggered", error);
    },
  });

  const handleLogin = async () => {
    if (!user) {
      try {
        login();
      } catch (error) {
        console.log("Error calling login function:", error);
      }
    }
  };

  const getWalletAddress = async () => {
    let counter = 0;
    let wallets = JSON.parse(localStorage.getItem("privy:connections") || "[]");

    // 30 sec max wait time
    while (wallets.length === 0 || (wallets[0].walletClientType !== "privy" && counter < 30)) {
      // Wait for 2 seconds before checking again
      await new Promise((resolve) => setTimeout(resolve, 2000));
      counter++;
      wallets = JSON.parse(localStorage.getItem("privy:connections") || "[]");
    }

    if (wallets.length > 0) {
      const walletAddress = wallets[0].address;
      localStorage.setItem("walletAddress", walletAddress);
      dispatch(setWalletAddress(walletAddress));
      return walletAddress;
    }
  };

  return (
    <div className="container flex justify-center">
      <div className="p-20 text-center">
        <h2 className="mb-4 text-2xl font-bold uppercase">Log In to Get Started</h2>
        <Button onClick={handleLogin}>Log in</Button>
      </div>
    </div>
  );
};

export default PrivyLogin;
