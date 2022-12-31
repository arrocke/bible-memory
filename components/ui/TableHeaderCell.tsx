import { ReactNode } from "react"

export interface TableHeaderCellProps {
  children?: ReactNode
  scope: 'row' | 'col'
  rowSpan?: number
  colSpan?: number
}

export default function TableHeaderCell({ children, scope, ...props }: TableHeaderCellProps) {
  return <th
    {...props}
    scope={scope}
    className={`
      text-left
      ${
        scope === 'col'
          ? 'text-sm'
          : 'h-10'
      }
    `}
  >
    {children}
  </th>
}