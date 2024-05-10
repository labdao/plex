import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { fetchNFTData } from './asyncActions';

interface NFTState {
  data: any;
  loading: boolean;
  error: string | null;
}

const initialState: NFTState = {
  data: {},
  loading: false,
  error: null,
};

export const nftListSlice = createSlice({
  name: 'nftList',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchNFTData.pending, (state) => {
      state.loading = true;
    });
    builder.addCase(fetchNFTData.fulfilled, (state, action: PayloadAction<any>) => {
      state.data = action.payload;
      state.loading = false;
    });
    builder.addCase(fetchNFTData.rejected, (state, action) => {
      state.error = action.error.message;
      state.loading = false;
    });
  },
});

export default nftListSlice.reducer;