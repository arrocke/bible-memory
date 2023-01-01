import { InputHTMLAttributes } from "react"

export interface NumberInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  onChange(value?: number): void
}

export default function NumberInput({ onChange, className = '', ...props }: NumberInputProps) {
  return <input
    className={`${className} px-2 py-1 rounded border border-gray-400 shadow-inner focus:outline outline-yellow-500 focus:border-yellow-500 `}
    {...props}
    type="number"
    onChange={e => onChange(e.target.value === '' ? undefined : e.target.valueAsNumber)}
  />
}