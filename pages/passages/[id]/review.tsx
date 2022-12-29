import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/router'
import VerseTyper, { ProgressUpdate } from '../../../components/VerseTyper'
import styles from './review.module.css'
import Link from 'next/link'

interface Passage { 
  id: string
  reference: string
  reviewDate?: Date
  level: number
  text: string
}

interface PassageJSON {
  id: string
  reference: string
  reviewDate?: string
  level: number
  text: string
}

export default function ReviewPage() {
  const router = useRouter()
  const [passage, setPassage] = useState<Passage>()
  const id = router.query.id as string

  const continueLink = useRef<HTMLAnchorElement>(null)

  useEffect(() => {
    if (typeof id === 'string') {
      loadPassage(id)
    }
  }, [id])

  async function loadPassage(id: string) {
    const request = await fetch(`/api/passages/${id}`)
    const body = await request.json() as PassageJSON
    setPassage(
      {
        ...body,
        reviewDate: body.reviewDate ? new Date(body.reviewDate) : undefined
      }
    )
  }

  async function createPassage(reference: string) {
    await fetch('api/passages', {
      method: 'POST',
      body: JSON.stringify({ reference }),
      headers: {
        'content-type': 'application/js'
      }
    })
  }

  const [progress, setProgress] = useState<ProgressUpdate>()
  const [isComplete, setComplete] = useState(false)

  useEffect(() => {
    const isComplete = progress ? progress.totalWords > 0 && progress.wordsComplete === progress.totalWords : false
    setComplete(isComplete)
    if (isComplete) {
      setTimeout(() => {
        continueLink.current?.focus()
      })
      // TODO: update persistent state
    }
  }, [progress])

  return (
    <div>
      <h1>Review {passage?.reference}</h1>
      <div>
        Progress: {progress ? (100 * progress.wordsComplete / progress.totalWords).toFixed(0) : 0}%
      </div>
      <div>
        Accuracy: {progress ? (100 * progress.correctWords / progress.totalWords).toFixed(0) : 0}%
      </div>
      {passage
        ?  <VerseTyper
            className={styles.typer}
            text={passage.text}
            onProgress={setProgress}
          />
        : null }
      <div style={{display: isComplete ? 'block' : 'none' }}>
        <Link ref={continueLink} href={'/passages'}>Continue</Link>
      </div>
    </div>
  )
}


