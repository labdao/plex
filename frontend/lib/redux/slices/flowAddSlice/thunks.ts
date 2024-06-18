import { setStripeCheckoutError, setStripeCheckoutLoading, setStripeCheckoutSuccess, setStripeCheckoutUrl } from '@/lib/redux'
import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { getCheckoutURL } from '../stripeCheckoutSlice/asyncActions'
import { createFlow } from './asyncActions'
import { Kwargs, setFlowAddError, setFlowAddID, setFlowAddLoading, setFlowAddSuccess } from './slice'

interface FlowPayload {
  name: string,
  toolCid: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addFlowThunk = createAppAsyncThunk(
  'flow/addFlow',
  async ({ name, toolCid, scatteringMethod, kwargs }: FlowPayload, { dispatch }) => {
    try {
      const response = await createFlow({ name, toolCid, scatteringMethod, kwargs })
      if (response && response.cid) {
        dispatch(setFlowAddSuccess(true))
        dispatch(setFlowAddID(response.ID))
      } else {
        console.log('Failed to add tool.', response)
        dispatch(setFlowAddError('Failed to add tool.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add Flow.', error)
      if (error instanceof Error) {
        dispatch(setFlowAddError(error.message))
      } else {
        dispatch(setFlowAddError('Failed to add tool.'))
      }
      return false
    }
  }
)

export const addFlowWithCheckoutThunk = createAppAsyncThunk(
  'flow/addFlowWithCheckout',
  async ({ toolCid, scatteringMethod, kwargs }: FlowPayload, { dispatch }) => {
    try {
      dispatch(setFlowAddError(null));
      dispatch(setStripeCheckoutError(null));
      dispatch(setStripeCheckoutLoading(true));

      const checkoutPayload = {
        toolCid,
        scatteringMethod,
        kwargs: JSON.stringify(kwargs),
      };
      const checkoutResponse = await getCheckoutURL(checkoutPayload);
      dispatch(setStripeCheckoutSuccess(true));
      dispatch(setStripeCheckoutUrl(checkoutResponse));
      dispatch(setStripeCheckoutLoading(false));

      return { checkout: checkoutResponse };
    } catch (error: unknown) {
      console.log('Failed to add flow with checkout.', error);
      if (error instanceof Error) {
        dispatch(setFlowAddError(error.message));
        dispatch(setStripeCheckoutError(error.message));
      } else {
        dispatch(setFlowAddError('Failed to add flow with checkout.'));
        dispatch(setStripeCheckoutError('Failed to add flow with checkout.'));
      }
      dispatch(setStripeCheckoutLoading(false));
      dispatch(setFlowAddLoading(false));
      return false;
    }
  }
);