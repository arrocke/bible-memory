import { ReactNode } from "react"

export interface TableHeaderProps {
  children: ReactNode
  className?: string
}

export default function TableHeader({ children, className = '' }: TableHeaderProps) {
  return <thead className={`${className} border-b border-black`}>
    {children}
  </thead>
}
