import Link from "../../components/ui/Link";
import { add, isSameDay, isBefore, startOfToday, format } from "date-fns";
import { useEffect, useState } from "react";
import Page from "../../components/ui/Page";
import PageTitle from "../../components/ui/PageTitle";
import Table from "../../components/ui/Table";
import TableHeader from "../../components/ui/TableHeader";
import TableHeaderCell from "../../components/ui/TableHeaderCell";
import TableDataCell from "../../components/ui/TableDataCell";
import TableBody from "../../components/ui/TableBody";
import TableFooter from "../../components/ui/TableFooter";
import PageHeader from "../../components/ui/PageHeader";

interface Passage {
  id: string;
  reference: string;
  reviewDate?: Date;
  level: number;
}

interface PassageJSON {
  id: string;
  reference: string;
  reviewDate?: string;
  level: number;
}

export default function Home() {
  const [passages, setPassages] = useState<Passage[]>([]);

  useEffect(() => {
    loadPassages();
  }, []);

  async function loadPassages() {
    const request = await fetch("/api/passages");
    const body = (await request.json()) as PassageJSON[];
    setPassages(
      body.map((passage) => {
        return {
          ...passage,
          reviewDate: passage.reviewDate
            ? add(new Date(passage.reviewDate), {
                minutes: new Date(passage.reviewDate).getTimezoneOffset(),
              })
            : undefined,
        };
      })
    );
  }

  return (
    <Page>
      <PageHeader>
        <PageTitle>Passages</PageTitle>
      </PageHeader>
      <div className="w-full overflow-x-auto sm:flex sm:justify-center whitespace-nowrap">
        <Table>
          <TableHeader>
            <tr>
              <TableHeaderCell className="px-2" scope="col">
                ID
              </TableHeaderCell>
              <TableHeaderCell className="px-2" scope="col">
                PASSAGE
              </TableHeaderCell>
              <TableHeaderCell className="px-2" scope="col">
                LEVEL
              </TableHeaderCell>
              <TableHeaderCell className="px-2" scope="col">
                NEXT REVIEW
              </TableHeaderCell>
              <TableHeaderCell className="px-2" scope="col"></TableHeaderCell>
            </tr>
          </TableHeader>
          <TableBody>
            {passages.map((passage) => {
              const today = startOfToday();
              const isDue = passage.reviewDate
                ? isSameDay(passage.reviewDate, today) ||
                  isBefore(passage.reviewDate, today)
                : false;
              return (
                <tr className={isDue ? "bg-green-200" : ""} key={passage.id}>
                  <TableDataCell>{passage.id}</TableDataCell>
                  <TableHeaderCell className="px-2" scope="row">
                    {passage.reference}
                  </TableHeaderCell>
                  <TableDataCell className="px-2">
                    {passage.level ?? "-"}
                  </TableDataCell>
                  <TableDataCell className="px-2">
                    {passage.reviewDate
                      ? format(passage.reviewDate, "MM/dd/yyyy")
                      : "-"}
                  </TableDataCell>
                  <TableDataCell className="px-2">
                    <Link className="mr-1" href={`/passages/${passage.id}`}>
                      Edit
                    </Link>
                    |
                    <Link
                      className="mx-1"
                      href={`/passages/${passage.id}/review?mode=learn`}
                    >
                      Learn
                    </Link>
                    |
                    <Link
                      className="mx-1"
                      href={`/passages/${passage.id}/review?mode=recall`}
                    >
                      Recall
                    </Link>
                    |
                    <Link
                      className="ml-1"
                      href={`/passages/${passage.id}/review?mode=review`}
                    >
                      Review
                    </Link>
                  </TableDataCell>
                </tr>
              );
            })}
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
  );
}
