import Link from '../../components/ui/Link'
import { format } from 'date-fns'
import { useEffect, useState } from 'react'
import Page from '../../components/ui/Page'
import PageTitle from '../../components/ui/PageTitle'
import Table from '../../components/ui/Table'
import TableHeader from '../../components/ui/TableHeader'
import TableHeaderCell from '../../components/ui/TableHeaderCell'
import TableDataCell from '../../components/ui/TableDataCell'
import TableBody from '../../components/ui/TableBody'
import TableFooter from '../../components/ui/TableFooter'
import PageHeader from '../../components/ui/PageHeader'
import EditableNumber from '../../components/ui/EditableNumber'

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

  async function updatePassageLevel({ id, level }: { id: string, level: number }) {
    await fetch(`/api/passages/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ level }),
      headers: {
        'content-type': 'application/js'
      }
    })
    await loadPassages()
  }

  return (
    <Page>
      <PageHeader>
        <PageTitle>Passages</PageTitle>
      </PageHeader>
      <Table className="w-full">
        <TableHeader>
          <tr>
            <TableHeaderCell scope="col">PASSAGE</TableHeaderCell>
            <TableHeaderCell scope="col">LEVEL</TableHeaderCell>
            <TableHeaderCell scope="col">NEXT REVIEW</TableHeaderCell>
            <TableHeaderCell scope="col"></TableHeaderCell>
          </tr>
        </TableHeader>
        <TableBody>
          {
            passages.map(passage => <tr key={passage.id}>
              <TableHeaderCell scope="row">{passage.reference}</TableHeaderCell>
              <TableDataCell>
                <EditableNumber
                  className="w-40"
                  value={passage.level}
                  onChange={(level) => updatePassageLevel({ id: passage.id, level })}
                  min={0}
                  max={9}
                />
              </TableDataCell>
              <TableDataCell>{passage.reviewDate ? format(passage.reviewDate, 'MM/dd/yyyy') : null}</TableDataCell>
              <TableDataCell>
                <Link className="mr-1" href={`/passages/${passage.id}/review?mode=learn`}>Learn</Link>
                |
                <Link className="mx-1" href={`/passages/${passage.id}/review?mode=recall`}>Recall</Link>
                |
                <Link className="ml-1" href={`/passages/${passage.id}/review?mode=review`}>Review</Link>
              </TableDataCell>
            </tr>)
          }
        </TableBody>
        <TableFooter>
          <tr>
            <td colSpan={4}>
              <Link href="/passages/new" className="py-1 block">
                + Add Passage
              </Link>
            </td>
          </tr>
        </TableFooter>
      </Table>
    </Page>
  )
}

