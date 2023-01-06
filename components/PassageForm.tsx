import { useState } from "react"
import Button from "./ui/Button"
import Label from "./ui/Label"
import MultilineTextInput from "./ui/MultilineTextInput"
import NumberInput from "./ui/NumberInput"
import TextInput from "./ui/TextInput"

interface Reference {
  book?: string
  startChapter?: number
  startVerse?: number
  endChapter?: number
  endVerse?: number
}

export interface PassageFormData {
  reference: Reference
  text: string
}

export interface PassageFormProps {
  initialData? : { text: string, reference: string }
  onSubmit(data: { text: string, reference: string }): void
  onDelete?(): void
}

const REFERENCE_MATCH = /^((?:(?:1|2) )?\w+) (\d+):(\d+)-(\d+):(\d+)$/
function parseReference(reference: string): Reference {
  const match = REFERENCE_MATCH.exec(reference)
  if (match) {
    return {
      book: match[1],
      startChapter: parseInt(match[2]),
      startVerse: parseInt(match[3]),
      endChapter: parseInt(match[4]),
      endVerse: parseInt(match[5])
    }
  } else {
    return {}
  }
}

export default function PassageForm({ initialData, onSubmit, onDelete }: PassageFormProps) {
  const [{ reference, text }, setData] = useState<PassageFormData>(() => ({
    reference: initialData?.reference ? parseReference(initialData.reference) : {},
    text: initialData?.text ?? ''
  }))

  return <form
    onSubmit={(e) => {
      e.preventDefault()
      onSubmit({
        reference: `${reference.book} ${reference.startChapter}:${reference.startVerse}-${reference.endChapter}:${reference.endVerse}`,
        text
      })
    }}
  >
    <fieldset className="flex">
      <Label as="legend">REFERENCE</Label>
      <div className="mr-6">
        <Label secondary>BOOK</Label>
        <TextInput 
          className="block w-32"
          id="reference"
          value={reference.book}
          onChange={book => setData(data => ({ ...data, reference: { ...reference, book } }))}
          autoComplete="off"
        />
      </div>
      <div className="mr-2">
        <Label secondary>Chapter</Label>
        <NumberInput
          className="block w-12"
          id="reference"
          type="number"
          value={reference.startChapter}
          onChange={startChapter => setData(data => ({ ...data, reference: { ...reference, startChapter } }))}
          autoComplete="off"
          step={1}
          min={1}
        />
      </div>
      <div className="pt-7 mr-2">
        :
      </div>
      <div className="mr-4">
        <Label secondary>Verse</Label>
        <NumberInput
          className="block w-12"
          id="reference"
          type="number"
          value={reference.startVerse}
          onChange={startVerse => setData(data => ({ ...data, reference: { ...reference, startVerse } }))}
          autoComplete="off"
          step={1}
          min={1}
        />
      </div>
      <div className="pt-7 mr-4">
        &mdash;
      </div>
      <div className="mr-2">
        <Label secondary>Chapter</Label>
        <NumberInput
          className="block w-12"
          id="reference"
          type="number"
          value={reference.endChapter}
          onChange={endChapter => setData(data => ({ ...data, reference: { ...reference, endChapter } }))}
          autoComplete="off"
          step={1}
          min={1}
        />
      </div>
      <div className="pt-7 mr-2">
        :
      </div>
      <div>
        <Label secondary>Verse</Label>
        <NumberInput
          className="block w-12"
          id="reference"
          type="number"
          value={reference.endVerse}
          onChange={endVerse => setData(data => ({ ...data, reference: { ...reference, endVerse } }))}
          autoComplete="off"
          step={1}
          min={1}
        />
      </div>
    </fieldset>
    <div className="mt-2">
      <Label htmlFor="text">TEXT</Label>
      <MultilineTextInput 
        className="block w-full h-96 resize-none"
        id="text"
        value={text}
        onChange={text => setData(data => ({ ...data, text }))}
      />
    </div>
    <div className="mt-4">
      <Button type="submit">Save</Button>
      {onDelete ? <Button className="ml-4" destructive onClick={onDelete}>Delete</Button> : null}
    </div>
  </form>
}