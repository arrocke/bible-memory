import { ReactNode } from "react"

export interface PageHeaderProps {
  children: ReactNode
}

export default function PageHeader({ children }: PageHeaderProps) {
  return <div className="mb-4 pt-2">
    {children}
  </div>
}

