import EditableField from "./EditableField"
import { format, parse } from 'date-fns'
import TextInput from "./TextInput"

export interface EditableDateProps {
  className?: string
  value?: Date
  onChange(value?: Date): void
}

console.log(parse('12/22/2022', 'MM/dd/yyyy', new Date()))

export default function EditableDate({ className = '', value, onChange }: EditableDateProps) {
  return <EditableField 
    className={`${className} w-44`}
    value={value ? format(value, 'MM/dd/yyyy') : ''}
    edit={({ value, onChange }) => 
      <TextInput
        className="flex-1 w-0"
        ref={(input) => input?.focus()}
        pattern="\d{2}\/\d{2}\/\d{4}"
        value={value}
        onChange={(value) => onChange(value)}
        placeholder="mm/dd/yyyy"
      />
    }
    display={<span className="w-24 inline-block">{value ? format(value, 'MM/dd/yyyy') : '-'}</span>}
    onSubmit={value => onChange(value ? parse(value, 'MM/dd/yyyy', new Date()) : undefined)}
  />
}