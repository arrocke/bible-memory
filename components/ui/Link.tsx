import BaseLink from 'next/link'
import { ComponentProps, forwardRef } from 'react'

interface LinkProps extends ComponentProps<typeof BaseLink> {
  button?: boolean
}

const Link = forwardRef<HTMLAnchorElement, LinkProps>(({ className = '', button, ...props }, ref) => {
  return <BaseLink {...props} ref={ref} className={`${className} font-bold focus:outline outline-yellow-500 outline-2 rounded ${button ? 'inline-block bg-blue-600 text-white px-3 py-1' : 'text-blue-600'}`} />
})

export default Link