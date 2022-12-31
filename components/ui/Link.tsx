import BaseLink from 'next/link'
import { ComponentProps } from 'react'

interface LinkProps extends ComponentProps<typeof BaseLink> {
  button?: boolean
}

export default function Link({ className = '', button, ...props }: LinkProps) {
  return <BaseLink {...props} className={`${className} font-bold focus:outline outline-yellow-500 outline-2 rounded ${button ? 'inline-block bg-blue-600 text-white px-3 py-1' : 'text-blue-600'}`} />
}