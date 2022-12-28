import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

interface Passage { 
  id: string
  reference: string
  reviewDate?: Date
  level: number
}

interface PassageJSON {
  id: string
  reference: string
  reviewDate?: string
  level: number
}

export default function ReviewPage() {
  const router = useRouter()
  const [passage, setPassage] = useState<Passage>()
  const id = router.query.id as string
  console.log(id)

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

  return (
    <div>
      <h1>Review {passage?.reference}</h1>
    </div>
  )
}


