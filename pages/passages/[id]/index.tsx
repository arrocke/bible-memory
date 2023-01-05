import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/router'
import VerseTyper, { ProgressUpdate, VerseTyperProps } from '../../../components/VerseTyper'
import Link from '../../../components/ui/Link'
import Page from '../../../components/ui/Page'
import PageHeader from '../../../components/ui/PageHeader'
import PageTitle from '../../../components/ui/PageTitle'
import BackLink from '../../../components/ui/BackLink'
import PassageForm from '../../../components/PassageForm'

interface Passage { 
  id: string
  reference: string
  text: string
}

export default function EditPassagePage() {
  const router = useRouter()
  const [passage, setPassage] = useState<Passage>()
  const id = router.query.id as string

  useEffect(() => {
    if (typeof id === 'string') {
      loadPassage(id)
    }
  }, [id])

  async function loadPassage(id: string) {
    const request = await fetch(`/api/passages/${id}`)
    setPassage(await request.json())
  }

  async function onSubmit(data: { text: string, reference: string }) {
    await fetch(`/api/passages/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
      headers: {
        'content-type': 'application/js'
      }
    })
    router.push('/passages')
  }

  return (
    <Page>
      <PageHeader>
        <BackLink href="/passages">Back to Passages</BackLink>
        <PageTitle>Edit {passage?.reference}</PageTitle>
      </PageHeader>
      {passage && <PassageForm initialData={passage} onSubmit={onSubmit}/>}
    </Page>
  )
}



