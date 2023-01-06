import { ButtonHTMLAttributes } from "react";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  destructive?: boolean
}

export default function Button({ className = '', type = 'button', destructive, ...props }: ButtonProps) {
  return <button
    {...props}
    type={type}
    className={`
      ${className}
      rounded font-bold px-3 py-1 focus:outline outline-yellow-500 outline-2
      ${destructive ? 'bg-red-600 text-white' : 'bg-blue-600 text-white'}
    `}
  />
}