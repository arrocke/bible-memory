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
              <TableDataCell>{passage.level}</TableDataCell>
              <TableDataCell>{passage.reviewDate ? format(passage.reviewDate, 'MM/dd/yyyy') : null}</TableDataCell>
              <TableDataCell>
                <Link href={`/passages/${passage.id}/review`}>Review</Link>
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

