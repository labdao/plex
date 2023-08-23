'use client'

/* Core */
import Link from 'next/link'

/* Instruments */
import styles from './topnav.module.css'

export const TopNav = () => {
  return (
    <nav className={styles.navbar}>
      <span className={styles.link}>Plex</span>
      {/* Other links or elements can be added here if required */}
    </nav>
  )
}
