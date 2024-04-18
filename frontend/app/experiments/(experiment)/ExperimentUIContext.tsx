import { createContext, useState } from "react";

interface ExperimentUIContextType {
  activeJobUUID: string | undefined;
  setActiveJobUUID: React.Dispatch<React.SetStateAction<string | undefined>>;
  activeCheckpointUrl: string | undefined;
  setActiveCheckpointUrl: React.Dispatch<React.SetStateAction<string | undefined>>;
  modelPanelOpen: boolean;
  setModelPanelOpen: React.Dispatch<React.SetStateAction<boolean>>;
}
export const ExperimentUIContext = createContext<ExperimentUIContextType>({} as ExperimentUIContextType);
export function ExperimentUIContextProvider({ children }: { children: React.ReactNode }) {
  const [activeJobUUID, setActiveJobUUID] = useState<string | undefined>(undefined);
  const [activeCheckpointUrl, setActiveCheckpointUrl] = useState<string | undefined>(undefined);
  const [modelPanelOpen, setModelPanelOpen] = useState<boolean>(false);

  return (
    <ExperimentUIContext.Provider
      value={{
        activeJobUUID,
        setActiveJobUUID,
        activeCheckpointUrl,
        setActiveCheckpointUrl,
        modelPanelOpen,
        setModelPanelOpen,
      }}
    >
      {children}
    </ExperimentUIContext.Provider>
  );
}
