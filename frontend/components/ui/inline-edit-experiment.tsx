import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { PencilIcon, CheckIcon } from "lucide-react";
import { AppDispatch, experimentUpdateThunk } from '@/lib/redux';

interface Experiment {
    ID: number;
    Name: string;
}

interface InlineEditExperimentProps {
    experiment: Experiment;
}

export const InlineEditExperiment: React.FC<InlineEditExperimentProps> = ({ experiment }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [name, setName] = useState(experiment.Name);
  const dispatch: AppDispatch = useDispatch();

  const handleRename = () => {
    const action = experimentUpdateThunk({
        experimentId: experiment.ID.toString(),
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
    <div className="flex items-center">
      <span className="flex-1 mr-2 truncate">{name}</span>
      <button onClick={() => setIsEditing(true)} className="absolute hidden group-hover:block right-[2px]">
        <PencilIcon size={16} />
      </button>
    </div>
  );
};

