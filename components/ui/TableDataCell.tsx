import { ReactNode } from "react"

export interface TableDataCellProps {
  children?: ReactNode
  rowSpan?: number
  colSpan?: number
  className?: string
}

export default function TableDataCell({ children, className, ...props }: TableDataCellProps) {
  return <td {...props} className={`h-8 ${className}`}>
    {children}
  </td>
}
