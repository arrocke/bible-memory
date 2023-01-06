import { FormEvent, ReactNode, useEffect, useState } from "react"
import Icon from "./Icon"

export interface EditableFieldProps<T> {
  className?: string
  value?: T
  edit: (props: { value?: T, onChange(value?: T): void }) => ReactNode
  display: ReactNode
  onSubmit(value?: T): void
  onCancel?(): void
}

export default function EditableField<T>({ className = '', value, edit, display, onSubmit, onCancel }: EditableFieldProps<T>) {
  const [isEditing, setEditing] = useState(false)
  const [mutableValue, setValue] = useState(value)
  useEffect(() => setValue(value), [value])

  function submit(e: FormEvent) {
    e.preventDefault()
    setEditing(false)
    onSubmit(mutableValue)
  }
  
  function cancel() {
    setEditing(false)
    onCancel?.()
  }

  if (isEditing) {
    return <form className={`${className} flex align-middle`} onSubmit={submit}>
      {edit({ value: mutableValue, onChange: setValue })}
      <button className="px-2 py-1 w-8 focus:outline-2 outline-yellow-500">
        <span className="sr-only">Submit</span>
        <Icon icon="check" aria-hidden="true" />
      </button>
      <button type="button" className="px-2 py-1 w-8 focus:outline-2 outline-yellow-500" onClick={cancel}>
        <span className="sr-only">Cancel</span>
        <Icon icon="close" aria-hidden="true" />
      </button>
    </form>
  } else {
    return <div className={`${className} flex items-center`}>
      {display}
      <button type="button" className="px-2 py-1 focus:outline-2 outline-yellow-500" onClick={() => setEditing(true)}>
        <span className="sr-only">Edit</span>
        <Icon icon="pencil" aria-hidden="true" />
      </button>
    </div>
  }
}
