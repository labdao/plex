'use client'

import styles from './footbar.module.css'

export const FootBar = () => {
  const sapling = '🌱'

  return (
    <div className={styles.footbar}>
      {Array(200).fill(sapling).join(' ')}
    </div>
  )
}
