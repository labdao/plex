'use client'

/* Core */
import { Provider } from 'react-redux'

/* Instruments */
import { reduxStore } from '@/lib/redux'

import { Web3AuthProvider } from './Web3AuthProvider'

export const Providers = (props: React.PropsWithChildren) => {
  return (
    <Provider store={reduxStore}>
      <Web3AuthProvider>
        {props.children}
      </Web3AuthProvider>
    </Provider>
  )
}
