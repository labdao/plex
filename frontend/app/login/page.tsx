'use client'

import { useEffect } from 'react'
import LoginComponent from '../components/Web3AuthLogin/LoginComponent'
import { useSelector, selectIsLoggedIn} from '@/lib/redux'
import { useRouter } from 'next/navigation'

export default function LoginPage() {
  const isUserLoggedIn = useSelector(selectIsLoggedIn)
  const router = useRouter()

  useEffect(() => {
    if (isUserLoggedIn) {
      router.push('/')
    }
  }, [router, isUserLoggedIn])

  return (
    <div>
      <LoginComponent />
    </div>
  )
}