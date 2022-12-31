import { InputHTMLAttributes } from "react"

export interface TextInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  onChange(value: string): void
}

export default function TextInput({ onChange, className = '', ...props }: TextInputProps) {
  return <input
    className={`${className} px-2 py-1 rounded border border-gray-400 shadow-inner focus:outline outline-yellow-500 focus:border-yellow-500 `}
    {...props}
    onChange={e => onChange(e.target.value)}
  />
}