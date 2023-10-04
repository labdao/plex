/* Instruments */
import {
  userSlice,
  dataFileAddSlice,
  dataFileListSlice,
  toolAddSlice,
  toolListSlice,
  flowAddSlice,
  flowListSlice,
} from './slices'

import { jobSlice } from './slices/jobSlice'

export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
  dataFileList: dataFileListSlice.reducer,
  toolAdd: toolAddSlice.reducer,
  toolList: toolListSlice.reducer,
  flowAdd: flowAddSlice.reducer,
  flowList: flowListSlice.reducer,
  job: jobSlice.reducer,
}
