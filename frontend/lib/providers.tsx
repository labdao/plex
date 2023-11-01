'use client'

import { Provider } from 'react-redux'

import { reduxStore } from '@/lib/redux'

import PrivyProviderWrapper from './PrivyProviderWrapper'

export const Providers = (props: React.PropsWithChildren) => {
  return (
    <Provider store={reduxStore}>
      <PrivyProviderWrapper>
        {props.children}
      </PrivyProviderWrapper>
    </Provider>
  )
}