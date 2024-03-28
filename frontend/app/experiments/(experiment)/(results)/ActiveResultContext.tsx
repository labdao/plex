import { createContext, useState } from "react";

interface ActiveResultContextType {
  activeJobUUID: string | undefined;
  setActiveJobUUID: React.Dispatch<React.SetStateAction<string | undefined>>;
  activeCheckpointUrl: string | undefined;
  setActiveCheckpointUrl: React.Dispatch<React.SetStateAction<string | undefined>>;
}
export const ActiveResultContext = createContext<ActiveResultContextType>({} as ActiveResultContextType);
export function ActiveResultContextProvider({ children }: { children: React.ReactNode }) {
  const [activeJobUUID, setActiveJobUUID] = useState<string | undefined>(undefined);
  const [activeCheckpointUrl, setActiveCheckpointUrl] = useState<string | undefined>(undefined);

  return (
    <ActiveResultContext.Provider
      value={{
        activeJobUUID,
        setActiveJobUUID,
        activeCheckpointUrl,
        setActiveCheckpointUrl,
      }}
    >
      {children}
    </ActiveResultContext.Provider>
  );
}
