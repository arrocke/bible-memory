import { ReactNode } from "react"

export interface TableDataCellProps {
  children?: ReactNode
  rowSpan?: number
  colSpan?: number
}

export default function TableDataCell({ children, ...props }: TableDataCellProps) {
  return <td {...props} className="h-8">
    {children}
  </td>
}
