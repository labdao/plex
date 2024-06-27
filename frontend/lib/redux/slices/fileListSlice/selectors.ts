import type { ReduxState } from '@/lib/redux'

export const selectFileList = (state: ReduxState) => state.fileList.files
export const selectFileListPagination = (state: ReduxState) => state.fileList.pagination
export const selectFileListLoading = (state: ReduxState) => state.fileList.status === 'loading';
export const selectFileListSuccess = (state: ReduxState) => state.fileList.success
export const selectFileListError = (state: ReduxState) => state.fileList.error
