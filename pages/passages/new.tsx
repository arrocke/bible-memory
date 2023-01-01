import Link from "../../components/ui/Link";
import { useRouter } from "next/router";
import { FormEvent, useState } from "react";
import Page from "../../components/ui/Page";
import PageHeader from "../../components/ui/PageHeader";
import PageTitle from "../../components/ui/PageTitle";
import BackLink from "../../components/ui/BackLink";
import TextInput from "../../components/ui/TextInput";
import MultilineTextInput from "../../components/ui/MultilineTextInput";
import Label from "../../components/ui/Label";
import Button from "../../components/ui/Button";
import NumberInput from "../../components/ui/NumberInput";

interface Reference {
  book?: string
  startChapter?: number
  startVerse?: number
  endChapter?: number
  endVerse?: number
}

export default function NewPassagePage() {
  const router = useRouter()
  const [reference, setReference] = useState<Reference>({})
  const [text, setText] = useState("");

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const referenceStr = `${reference.book} ${reference.startChapter}:${reference.startVerse}-${reference.endChapter}:${reference.endVerse}`
    await fetch("/api/passages", {
      method: "POST",
      body: JSON.stringify({ reference: referenceStr, text }),
      headers: {
        "content-type": "application/js",
      },
    });
    router.push('/passages')
  }

  return (
    <Page>
      <PageHeader>
        <BackLink href="/passages">Back to Passages</BackLink>
        <PageTitle>New Passage</PageTitle>
      </PageHeader>
      <form onSubmit={onSubmit}>
        <fieldset className="flex">
          <Label as="legend">REFERENCE</Label>
          <div className="mr-6">
            <Label secondary>BOOK</Label>
            <TextInput 
              className="block w-32"
              id="reference"
              value={reference.book}
              onChange={book => setReference(reference => ({ ...reference, book }))}
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
              onChange={startChapter => setReference(reference => ({ ...reference, startChapter }))}
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
              onChange={startVerse => setReference(reference => ({ ...reference, startVerse }))}
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
              onChange={endChapter => setReference(reference => ({ ...reference, endChapter }))}
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
              onChange={endVerse => setReference(reference => ({ ...reference, endVerse }))}
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
            onChange={setText}
          />
        </div>
        <div className="mt-4">
          <Button type="submit">Save</Button>
        </div>
      </form>
    </Page>
  );
}
