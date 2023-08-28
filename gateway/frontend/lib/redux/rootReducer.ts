/* Instruments */
import { dataFileAddSlice, userSlice } from './slices'
import { toolAddSlice } from './slices/toolAddSlice'

export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
  toolAdd: toolAddSlice.reducer,
}
