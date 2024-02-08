/* Instruments */
import {
  apiKeyAddSlice,
  apiKeyListSlice,
  dataFileAddSlice,
  dataFileListSlice,
  flowAddSlice,
  flowDetailSlice,
  flowListSlice,
  jobDetailSlice,
  toolAddSlice,
  toolDetailSlice,
  toolListSlice,
  userSlice,
} from "./slices";

export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
  dataFileList: dataFileListSlice.reducer,
  toolAdd: toolAddSlice.reducer,
  toolList: toolListSlice.reducer,
  toolDetail: toolDetailSlice.reducer,
  flowAdd: flowAddSlice.reducer,
  flowList: flowListSlice.reducer,
  flowDetail: flowDetailSlice.reducer,
  jobDetail: jobDetailSlice.reducer,
  apiKeyAdd: apiKeyAddSlice.reducer,
  apiKeyList: apiKeyListSlice.reducer,
};
