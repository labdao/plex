import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { PencilIcon, CheckIcon } from "lucide-react";
import { AppDispatch, flowUpdateThunk } from '@/lib/redux';

interface Flow {
    ID: number;
    Name: string;
}

interface InlineEditExperimentProps {
    flow: Flow;
}

export const InlineEditExperiment: React.FC<InlineEditExperimentProps> = ({ flow }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [name, setName] = useState(flow.Name);
  const dispatch: AppDispatch = useDispatch();

  const handleRename = () => {
    const action = flowUpdateThunk({
        flowId: flow.ID.toString(),
        updates: { name }
      });
      dispatch(action);
      setIsEditing(false);
  };

  if (isEditing) {
    return (
      <div className="flex items-center space-x-2">
        <input
          type="text"
          className="text-sm px-2 py-1 rounded"
          value={name}
          onChange={(e) => setName(e.target.value)}
          onBlur={handleRename}
        />
        <button onClick={handleRename} className="text-green-500">
          <CheckIcon size={16} />
        </button>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-between">
      <span>{name}</span>
      <button onClick={() => setIsEditing(true)} className="text-gray-500">
        <PencilIcon size={16} />
      </button>
    </div>
  );
};

