import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/router'
import VerseTyper, { ProgressUpdate, VerseTyperProps } from '../../../components/VerseTyper'
import Link from '../../../components/ui/Link'
import Page from '../../../components/ui/Page'
import PageHeader from '../../../components/ui/PageHeader'
import PageTitle from '../../../components/ui/PageTitle'
import BackLink from '../../../components/ui/BackLink'

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
  const mode = router.query.mode as VerseTyperProps['mode']

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

  async function updatePassage(id: string, review: boolean) {
    await fetch(`/api/passages/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ review }),
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
      if (mode === 'review') {
        updatePassage(id, progress?.totalWords === progress?.correctWordsWithHelp)
      }
    }
  }, [progress])

  return (
    <Page>
      <PageHeader>
        <BackLink href="/passages">Back to Passages</BackLink>
        <PageTitle>Review {passage?.reference}</PageTitle>
      </PageHeader>
      <div>
        Progress: {progress ? (100 * progress.wordsComplete / progress.totalWords).toFixed(0) : 0}%
      </div>
      <div className="mb-2">
        Accuracy: {progress ? (100 * progress.correctWords / progress.totalWords).toFixed(0) : 0}%
      </div>
      {passage
        ?  <VerseTyper
            className="mb-4"
            text={passage.text}
            mode={mode}
            onProgress={setProgress}
          />
        : null }
      <div style={{display: isComplete ? 'block' : 'none' }}>
        <Link button ref={continueLink} href={'/passages'}>Continue</Link>
      </div>
    </Page>
  )
}


