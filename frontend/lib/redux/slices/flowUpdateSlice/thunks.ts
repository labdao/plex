import { createAppAsyncThunk } from 'lib/redux/createAppAsyncThunk';
import { AppDispatch } from 'lib/redux/store';

import { updateFlow } from './asyncActions';
import { setFlowUpdateError, setFlowUpdateLoading, setFlowUpdateSuccess } from './slice';

interface UpdateFlowArgs {
    flowId: string;
}

export const flowUpdateThunk = createAppAsyncThunk(
    'flow/updateFlow',
    async ({ flowId }: UpdateFlowArgs, { dispatch }: { dispatch: AppDispatch }) => {
        dispatch(setFlowUpdateLoading(true));
        try {
            const result = await updateFlow(flowId);
            dispatch(setFlowUpdateSuccess(true));
            return result;
        } catch (error) {
            dispatch(setFlowUpdateError(error instanceof Error ? error.toString() : 'Failed to update Flow.'));
            return { error: error instanceof Error ? error.toString() : 'Failed to update Flow.' };
        } finally {
            dispatch(setFlowUpdateLoading(false));
        }
    }
);