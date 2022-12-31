import { ReactNode } from "react"

export interface PageProps {
  children: ReactNode
}

export default function Page({ children }: PageProps) {
  return <div className="lg:container mx-auto px-4">
    {children}
  </div>
}