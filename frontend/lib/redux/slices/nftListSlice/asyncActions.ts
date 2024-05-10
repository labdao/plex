import { createAsyncThunk } from '@reduxjs/toolkit';

import { getNFTContract } from '../../../ethereumConnection';

export const fetchNFTData = createAsyncThunk('nft/fetchData', async (tokenId: number) => {
    const nftContract = getNFTContract();
    const data = await nftContract.fetchData(tokenId);
    return data;
});