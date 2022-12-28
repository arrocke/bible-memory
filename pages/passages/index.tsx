import Link from 'next/link'
import { format } from 'date-fns'
import { useEffect, useState } from 'react'

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

export default function Home() {
  const [passages, setPassages] = useState<Passage[]>([])

  useEffect(() => {
    loadPassages()
  }, [])

  async function loadPassages() {
    const request = await fetch('/api/passages')
    const body = await request.json() as PassageJSON[]
    setPassages(body.map(passage => (
      {
        ...passage,
        reviewDate: passage.reviewDate ? new Date(passage.reviewDate) : undefined
      }
    )))
  }

  async function createPassage(reference: string) {
    await fetch('api/passages', {
      method: 'POST',
      body: JSON.stringify({ reference }),
      headers: {
        'content-type': 'application/js'
      }
    })
    await loadPassages()
  }

  return (
    <div>
      <h1>Passages</h1>
      <table>
        <thead>
          <tr>
            <th>Passage</th>
            <th>Level</th>
            <th>Next Review</th>
            <th></th>
          </tr>
          <tr>
            <td rowSpan={4}><button type="button" onClick={() => createPassage('Psalm 3')}>+ Add Passage</button></td>
          </tr>
        </thead>
        <tbody>
          {
            passages.map(passage => <tr key={passage.id}>
              <td>{passage.reference}</td>
              <td>{passage.level}</td>
              <td>{passage.reviewDate ? format(passage.reviewDate, 'MM/dd/yyyy') : null}</td>
              <td><Link href={`/passages/${passage.id}/review`}>Review</Link></td>
            </tr>)
          }
        </tbody>
      </table>
    </div>
  )
}

