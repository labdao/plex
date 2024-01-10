import type { ReduxState } from '@/lib/redux'

export const selectDataFileList = (state: ReduxState) => state.dataFileList.dataFiles
export const selectDataFileListPagination = (state: ReduxState) => state.dataFileList.pagination
export const selectDataFileListLoading = (state: ReduxState) => state.dataFileList.status === 'loading';
export const selectDataFileListSuccess = (state: ReduxState) => state.dataFileList.success
export const selectDataFileListError = (state: ReduxState) => state.dataFileList.error
