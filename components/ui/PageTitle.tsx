import { ReactNode } from "react"

export interface PageProps {
  children: ReactNode
}

export default function PageTitle({ children }: PageProps) {
  return <h1 className="text-xl font-bold">
    {children}
  </h1>
}
