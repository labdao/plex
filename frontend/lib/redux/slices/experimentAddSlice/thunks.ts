import { setStripeCheckoutError, setStripeCheckoutLoading, setStripeCheckoutSuccess, setStripeCheckoutUrl, selectUserSubscriptionStatus } from '@/lib/redux'
import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { getCheckoutURL } from '../stripeCheckoutSlice/asyncActions'
import { createExperiment } from './asyncActions'
import { Kwargs, setExperimentAddError, setExperimentAddID, setExperimentAddLoading, setExperimentAddSuccess } from './slice'

interface ExperimentPayload {
  name: string,
  modelId: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addExperimentThunk = createAppAsyncThunk(
  'experiment/addExperiment',
  async ({ name, modelId, scatteringMethod, kwargs }: ExperimentPayload, { dispatch }) => {
    try {
      const response = await createExperiment({ name, modelId, scatteringMethod, kwargs })
      if (response && response.id) {
        dispatch(setExperimentAddSuccess(true))
        dispatch(setExperimentAddID(response.ID))
      } else {
        console.log('Failed to add model.', response)
        dispatch(setExperimentAddError('Failed to add model.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add Experiment.', error)
      if (error instanceof Error) {
        dispatch(setExperimentAddError(error.message))
      } else {
        dispatch(setExperimentAddError('Failed to add model.'))
      }
      return false
    }
  }
)

export const addExperimentWithCheckoutThunk = createAppAsyncThunk(
  'experiment/addExperimentWithCheckout',
  async ({ name, modelId, scatteringMethod, kwargs }: ExperimentPayload, { dispatch, getState }) => {
    try {
      dispatch(setExperimentAddError(null));
      dispatch(setStripeCheckoutError(null));
      
      const state = getState();
      const subscriptionStatus = selectUserSubscriptionStatus(state);

      if (subscriptionStatus === 'active') {
        // User is subscribed, directly add the experiment
        return dispatch(addExperimentThunk({ name, modelId, scatteringMethod, kwargs })).unwrap();
      } else {
        // User is not subscribed, initiate checkout process
        dispatch(setStripeCheckoutLoading(true));
        const checkoutResponse = await getCheckoutURL();
        dispatch(setStripeCheckoutSuccess(true));
        dispatch(setStripeCheckoutUrl(checkoutResponse));
        dispatch(setStripeCheckoutLoading(false));
        return { checkout: checkoutResponse };
      }
    } catch (error: unknown) {
      console.log('Failed to process experiment request.', error);
      if (error instanceof Error) {
        dispatch(setExperimentAddError(error.message));
        dispatch(setStripeCheckoutError(error.message));
      } else {
        dispatch(setExperimentAddError('Failed to process experiment request.'));
        dispatch(setStripeCheckoutError('Failed to process experiment request.'));
      }
      dispatch(setStripeCheckoutLoading(false));
      dispatch(setExperimentAddLoading(false));
      return false;
    }
  }
);