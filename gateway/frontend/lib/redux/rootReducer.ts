/* Instruments */
import { dataFileAddSlice, userSlice } from './slices'

export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
}
