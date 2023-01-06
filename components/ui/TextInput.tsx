import { forwardRef, InputHTMLAttributes } from "react"

export interface TextInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  onChange(value: string): void
}

const TextInput = forwardRef<HTMLInputElement, TextInputProps>(({ onChange, value, className = '', ...props }, ref) => {
  return <input
    ref={ref}
    className={`${className} px-2 py-1 rounded border border-gray-400 shadow-inner focus:outline outline-yellow-500 focus:border-yellow-500 `}
    {...props}
    value={value ?? ''}
    onChange={e => onChange(e.target.value)}
  />
})

export default TextInput