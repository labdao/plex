import { setStripeCheckoutError, setStripeCheckoutLoading, setStripeCheckoutSuccess, setStripeCheckoutUrl } from '@/lib/redux'
import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { getCheckoutURL } from '../stripeCheckoutSlice/asyncActions'
import { createExperiment } from './asyncActions'
import { Kwargs, setExperimentAddError, setExperimentAddID, setExperimentAddLoading, setExperimentAddSuccess } from './slice'

interface ExperimentPayload {
  name: string,
  modelCid: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addExperimentThunk = createAppAsyncThunk(
  'experiment/addExperiment',
  async ({ name, modelCid, scatteringMethod, kwargs }: ExperimentPayload, { dispatch }) => {
    try {
      const response = await createExperiment({ name, modelCid, scatteringMethod, kwargs })
      if (response && response.cid) {
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
  async ({ modelCid, scatteringMethod, kwargs }: ExperimentPayload, { dispatch }) => {
    try {
      dispatch(setExperimentAddError(null));
      dispatch(setStripeCheckoutError(null));
      dispatch(setStripeCheckoutLoading(true));

      const checkoutPayload = {
        modelCid,
        scatteringMethod,
        kwargs: JSON.stringify(kwargs),
      };
      const checkoutResponse = await getCheckoutURL(checkoutPayload);
      dispatch(setStripeCheckoutSuccess(true));
      dispatch(setStripeCheckoutUrl(checkoutResponse));
      dispatch(setStripeCheckoutLoading(false));

      return { checkout: checkoutResponse };
    } catch (error: unknown) {
      console.log('Failed to add experiment with checkout.', error);
      if (error instanceof Error) {
        dispatch(setExperimentAddError(error.message));
        dispatch(setStripeCheckoutError(error.message));
      } else {
        dispatch(setExperimentAddError('Failed to add experiment with checkout.'));
        dispatch(setStripeCheckoutError('Failed to add experiment with checkout.'));
      }
      dispatch(setStripeCheckoutLoading(false));
      dispatch(setExperimentAddLoading(false));
      return false;
    }
  }
);