import { LabelHTMLAttributes } from "react";

export interface LabelProps extends LabelHTMLAttributes<HTMLLabelElement> {
  secondary?: boolean
}

export default function Label({ className = '', secondary, ...props }: LabelProps) {
  return <label {...props} className={`${className} font-bold ${secondary ? 'text-xs' : 'text-sm'}`} />
}