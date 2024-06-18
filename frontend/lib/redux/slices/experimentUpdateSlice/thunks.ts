import { createAppAsyncThunk } from 'lib/redux/createAppAsyncThunk';
import { AppDispatch } from 'lib/redux/store';

import { updateExperiment } from './asyncActions';
import { setExperimentUpdateError, setExperimentUpdateLoading, setExperimentUpdateSuccess } from './slice';

interface UpdateExperimentArgs {
    experimentId: string;
    updates: {
        name?: string;
        public?: boolean;
    };
}

export const experimentUpdateThunk = createAppAsyncThunk(
    'experiment/updateExperiment',
    async ({ experimentId, updates }: UpdateExperimentArgs, { dispatch }: { dispatch: AppDispatch }) => {
        dispatch(setExperimentUpdateLoading(true));
        try {
            const result = await updateExperiment(experimentId, updates);
            dispatch(setExperimentUpdateSuccess(true));
            return result;
        } catch (error) {
            dispatch(setExperimentUpdateError(error instanceof Error ? error.toString() : 'Failed to update Experiment.'));
            return { error: error instanceof Error ? error.toString() : 'Failed to update Experiment.' };
        } finally {
            dispatch(setExperimentUpdateLoading(false));
        }
    }
);