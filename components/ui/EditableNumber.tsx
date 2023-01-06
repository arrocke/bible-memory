import EditableField from "./EditableField"
import NumberStepper from "./NumberStepper"

export interface EditableNumberProps {
  className?: string
  value: number
  onChange(value: number): void
  min: number
  max: number
}

export default function EditableNumber({ className, value, onChange, min, max }: EditableNumberProps) {
  return <EditableField 
    className={className}
    value={value}
    edit={({ value, onChange }) => 
      <NumberStepper
        ref={(input) => input?.focus()}
        className="flex-grow"
        value={value ?? 0}
        onChange={onChange}
        min={min}
        max={max}
      />
    }
    display={<span className="min-w-[16px]">{value}</span>}
    onSubmit={onChange}
  />
}