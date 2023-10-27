"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useRouter } from "next/navigation";
import { useContext, useEffect, useState } from "react";
import { ReactNode } from "react";

import { selectWalletAddress, setIsLoggedIn, setWalletAddress, useDispatch, useSelector } from "@/lib/redux";

type UserLoaderProps = {
  children: ReactNode;
};

const UserLoader = ({ children }: UserLoaderProps) => {
  const dispatch = useDispatch();
  const router = useRouter();
  const [isLoaded, setIsLoaded] = useState(false);
  const { ready, authenticated } = usePrivy();

  const walletAddressFromRedux = useSelector(selectWalletAddress);

  useEffect(() => {
    const walletAddressFromLocalStorage = localStorage.getItem("walletAddress");

    if (!walletAddressFromRedux && walletAddressFromLocalStorage) {
      dispatch(setWalletAddress(walletAddressFromLocalStorage));
    }

    if (ready) {
      if (!authenticated) {
        router.push("/login");
      } else {
        dispatch(setIsLoggedIn(true));
      }
    }
    setIsLoaded(true);
  }, [dispatch, ready, authenticated]);

  if (!isLoaded) return null;

  return children;
};

export default UserLoader;
