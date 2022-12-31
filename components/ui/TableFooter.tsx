import { ReactNode } from "react"

export interface TableFooterProps {
  children: ReactNode
  className?: string
}

export default function TableFooter({ children, className = '' }: TableFooterProps) {
  return <thead className={`${className} border-t border-black`}>
    {children}
  </thead>
}
