/* Instruments */
import { counterSlice, userSlice } from './slices'

export const reducer = {
  user: userSlice.reducer,
  counter: counterSlice.reducer,
}
