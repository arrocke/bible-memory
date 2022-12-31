import { ButtonHTMLAttributes } from "react";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {}

export default function Button({ className = '', type = 'button', ...props }: ButtonProps) {
  return <button {...props} type={type} className={`${className} rounded bg-blue-600 text-white font-bold px-3 py-1 focus:outline outline-yellow-500 outline-2`} />
}