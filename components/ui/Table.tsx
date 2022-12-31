import { ReactNode } from "react"

export interface TableProps {
  children: ReactNode
  className?: string
}

export default function Table({ children, className }: TableProps) {
  return <table className={className}>
    {children}
  </table>
}