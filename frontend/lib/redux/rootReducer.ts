/* Instruments */
import { dataFileAddSlice, userSlice, toolAddSlice } from './slices'
import { jobSlice } from './slices/jobSlice'

export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
  toolAdd: toolAddSlice.reducer,
  job: jobSlice.reducer,
}
