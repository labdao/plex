'use client'

import { useEffect } from 'react'
import PrivyLoginComponent from '../components/PrivyLogin/PrivyLoginComponent'
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
      <PrivyLoginComponent />
    </div>
  )
}