import { createElement, LabelHTMLAttributes } from "react";

export interface LabelProps extends LabelHTMLAttributes<HTMLLabelElement> {
  secondary?: boolean
  as?: 'legend' | 'label'
}

export default function Label({ as = 'label', className = '', secondary, ...props }: LabelProps) {
  return createElement(as, {
    ...props,
    className: `${className} font-bold ${secondary ? 'text-xs' : 'text-sm'}`
  })
}