import { useState, useEffect } from 'react'

const MOBILE_BREAKPOINT = 768

export function use_is_mobile(): boolean {
  const [is_mobile, set_is_mobile] = useState(() =>
    typeof window !== 'undefined' ? window.innerWidth < MOBILE_BREAKPOINT : false
  )

  useEffect(() => {
    const handle_resize = () => {
      set_is_mobile(window.innerWidth < MOBILE_BREAKPOINT)
    }

    window.addEventListener('resize', handle_resize)
    return () => window.removeEventListener('resize', handle_resize)
  }, [])

  return is_mobile
}
