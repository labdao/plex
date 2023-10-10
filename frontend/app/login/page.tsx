'use client'

import { useEffect } from 'react'
import PrivyLoginComponent from '../components/PrivyLogin/PrivyLoginComponent'
import { useSelector, selectIsLoggedIn} from '@/lib/redux'
import { usePrivy } from '@privy-io/react-auth';
import { useRouter } from 'next/navigation'

export default function LoginPage() {
  const isUserLoggedIn = useSelector(selectIsLoggedIn)
  const { ready, authenticated } = usePrivy();
  const router = useRouter()

  useEffect(() => {
    if (ready) {
      if (authenticated) {
        router.push('/');
      }
    }
  }, [router, isUserLoggedIn])

  return (
    <div>
      <PrivyLoginComponent />
    </div>
  )
}