import Link from './Link'
import { ComponentProps, ReactNode } from 'react'
import Icon from './Icon'

export interface BackLinkProps {
  children: ReactNode
  href: ComponentProps<typeof Link>['href']
}

export default function BackLink({ children, ...props }: BackLinkProps) {
  return <Link {...props} className={'mb-1 inline-block text-sm'}>
    <Icon icon="left-long" className="mr-1" aria-hidden="true"/>
    {children}
  </Link>
}
