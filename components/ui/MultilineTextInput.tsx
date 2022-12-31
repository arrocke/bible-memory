import { TextareaHTMLAttributes } from "react"

export interface MultilineTextInputProps extends Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, 'onChange'> {
  onChange(value: string): void
}

export default function MultilineTextInput({ onChange, className = '', ...props }: MultilineTextInputProps) {
  return <textarea
    className={`${className} px-2 py-1 rounded border border-gray-400 shadow-inner focus:outline outline-yellow-500 focus:border-yellow-500 `}
    {...props}
    onChange={e => onChange(e.target.value)}
  />
}