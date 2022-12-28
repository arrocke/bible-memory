import VerseTyper from '../components/VerseTyper'
import styles from '../styles/Home.module.css'

export default function Home() {
  return (
    <div className={styles.container}>
      <main className={styles.main}>
        <VerseTyper className={styles.input} text="To you, O Lord, I lift up my soul." />
      </main>
    </div>
  )
}
