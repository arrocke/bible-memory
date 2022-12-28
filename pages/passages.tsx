import styles from '../styles/passages.module.css'
import { format } from 'date-fns'
import { useEffect, useState } from 'react'

interface Passage { 
  reference: string
  reviewDate?: Date
}

interface PassageJSON {
  reference: string
  reviewDate?: string
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
            <th>Next Review</th>
            <th></th>
          </tr>
          <tr>
            <td rowSpan={3}><button type="button" onClick={() => createPassage('Psalm 3')}>+ Add Passage</button></td>
          </tr>
        </thead>
        <tbody>
          {
            passages.map(passage => <tr key={passage.reference}>
              <td>{passage.reference}</td>
              <td>{passage.reviewDate ? format(passage.reviewDate, 'MM/dd/yyyy') : null}</td>
              <td>Review</td>
            </tr>)
          }
        </tbody>
      </table>
    </div>
  )
}

