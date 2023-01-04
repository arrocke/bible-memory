import { forwardRef } from "react";
import Icon from "./Icon";

export interface NumberStepperProps {
  className?: string
  value: number
  onChange(value: number): void
  min: number
  max: number
}

export const NumberStepper = forwardRef<HTMLInputElement, NumberStepperProps>(({ value, onChange, min, max, className = '' }, ref) => {
  function emitIfValid(value: number) {
    if (value >= min && value <= max ){ 
      onChange(value)
    }
  }

  return <div className={`${className} inline-flex rounded border border-gray-400 focus-within:outline outline-yellow-500 focus-within:border-yellow-500`}>
    <button
      className="rounded-l border-r border-gray-400 flex-none px-2 py-1"
      type="button"
      aria-hidden="true"
      onClick={() => emitIfValid(value - 1)}
      tabIndex={-1}
    >
      <Icon icon="minus" />
    </button>
    <input
      ref={ref}
      className="w-0 flex-1 px-2 py-1 focus:outline-none"
      type="number"
      value={value}
      onChange={e => emitIfValid(e.target.valueAsNumber)}
    />
    <button
      className="rounded-r border-l border-gray-400 flex-none px-2 py-1"
      type="button"
      aria-hidden="true"
      onClick={() => emitIfValid(value + 1)}
      tabIndex={-1}
    >
      <Icon icon="plus" />
    </button>
  </div>
})

export default NumberStepper