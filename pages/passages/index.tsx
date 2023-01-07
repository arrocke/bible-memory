import Link from '../../components/ui/Link'
import { add, isSameDay, isBefore, startOfToday } from 'date-fns'
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
import NumberStepper from '../../components/ui/NumberStepper'
import EditableField from '../../components/ui/EditableField'
import EditableDate from '../../components/ui/EditableDate'

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
    setPassages(body.map(passage => {
      return {
        ...passage,
        reviewDate: passage.reviewDate
          ? add(new Date(passage.reviewDate), { minutes: new Date(passage.reviewDate).getTimezoneOffset() })
          : undefined
      }
    }))
  }

  async function updatePassage({ id, ...data }: { id: string, level?: number, reviewDate?: Date }) {
    await fetch(`/api/passages/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
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
      <div className="w-full overflow-x-auto">
        <Table className="w-full min-w-[700px]">
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
              passages.map(passage => {
                const today = startOfToday()
                const isDue = passage.reviewDate
                  ? isSameDay(passage.reviewDate, today) || isBefore(passage.reviewDate, today)
                  : false
                return <tr className={isDue ? 'bg-green-200' : ''} key={passage.id}>
                  <TableHeaderCell scope="row">{passage.reference}</TableHeaderCell>
                  <TableDataCell>
                    <EditableNumber
                      className="w-40"
                      value={passage.level}
                      onChange={(level) => updatePassage({ id: passage.id, level })}
                      min={0}
                      max={9}
                    />
                  </TableDataCell>
                  <TableDataCell>
                    <EditableDate 
                      value={passage.reviewDate}
                      onChange={(reviewDate) => updatePassage({ id: passage.id, reviewDate })}
                    />
                  </TableDataCell>
                  <TableDataCell>
                    <Link className="mr-1" href={`/passages/${passage.id}`}>Edit</Link>
                    |
                    <Link className="mx-1" href={`/passages/${passage.id}/review?mode=learn`}>Learn</Link>
                    |
                    <Link className="mx-1" href={`/passages/${passage.id}/review?mode=recall`}>Recall</Link>
                    |
                    <Link className="ml-1" href={`/passages/${passage.id}/review?mode=review`}>Review</Link>
                  </TableDataCell>
                </tr>
              })
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
      </div>
    </Page>
  )
}

