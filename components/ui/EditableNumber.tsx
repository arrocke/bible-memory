import { FormEvent, useEffect, useRef, useState } from "react"
import Icon from "./Icon"
import NumberStepper from "./NumberStepper"

export interface EditableNumberProps {
  className?: string
  value: number
  onChange(value: number): void
  min: number
  max: number
}

export default function EditableNumber({ className = '', value, onChange, min, max }: EditableNumberProps) {
  const [isEditing, setEditing] = useState(false)
  const [mutableValue, setValue] = useState(value)
  useEffect(() => setValue(value), [value])

  function submit(e: FormEvent) {
    e.preventDefault()
    setEditing(false)
    onChange(mutableValue)
  }
  
  function cancel() {
    setEditing(false)
    setValue(value)
  }

  return <div className={`${className}`}>
    {
      isEditing
        ? <form className="flex" onSubmit={submit}>
            <NumberStepper ref={(input) => input?.focus()} className="flex-grow" value={mutableValue} onChange={setValue} min={min} max={max} />
            <button className="px-2 py-1">
              <span className="sr-only">Submit</span>
              <Icon icon="check" aria-hidden="true" />
            </button>
            <button type="button" className="px-2 py-1" onClick={cancel}>
              <span className="sr-only">Cancel</span>
              <Icon icon="close" aria-hidden="true" />
            </button>
          </form>
        : <>
            {value}
            <button type="button" className="px-2 py-1" onClick={() => setEditing(true)}>
              <span className="sr-only">Edit</span>
              <Icon icon="pencil" aria-hidden="true" />
            </button>
          </>
    }
  </div>
}