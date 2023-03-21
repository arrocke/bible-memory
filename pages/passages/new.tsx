import { useRouter } from "next/router";
import Page from "../../components/ui/Page";
import PageHeader from "../../components/ui/PageHeader";
import PageTitle from "../../components/ui/PageTitle";
import BackLink from "../../components/ui/BackLink";
import Label from "../../components/ui/Label";
import { FormEvent, useState } from "react";
import Button from "../../components/ui/Button";
import MultilineTextInput from "../../components/ui/MultilineTextInput";
import NumberInput from "../../components/ui/NumberInput";
import TextInput from "../../components/ui/TextInput";

interface Reference {
  book?: string;
  startChapter?: number;
  startVerse?: number;
  endChapter?: number;
  endVerse?: number;
}

export interface PassageFormData {
  reference: Reference;
  text: string;
}

export default function NewPassagePage() {
  const router = useRouter();

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    await fetch("/api/passages", {
      method: "POST",
      body: JSON.stringify({
        reference: `${reference.book} ${reference.startChapter}:${reference.startVerse}-${reference.endChapter}:${reference.endVerse}`,
        text,
      }),
      headers: {
        "content-type": "application/js",
      },
    });
    router.push("/passages");
  }

  const [{ reference, text }, setData] = useState<PassageFormData>(() => ({
    reference: {},
    text: "",
  }));

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
              onChange={(book) =>
                setData((data) => ({
                  ...data,
                  reference: { ...reference, book },
                }))
              }
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
              onChange={(startChapter) =>
                setData((data) => ({
                  ...data,
                  reference: { ...reference, startChapter },
                }))
              }
              autoComplete="off"
              step={1}
              min={1}
            />
          </div>
          <div className="pt-7 mr-2">:</div>
          <div className="mr-4">
            <Label secondary>Verse</Label>
            <NumberInput
              className="block w-12"
              id="reference"
              type="number"
              value={reference.startVerse}
              onChange={(startVerse) =>
                setData((data) => ({
                  ...data,
                  reference: { ...reference, startVerse },
                }))
              }
              autoComplete="off"
              step={1}
              min={1}
            />
          </div>
          <div className="pt-7 mr-4">&mdash;</div>
          <div className="mr-2">
            <Label secondary>Chapter</Label>
            <NumberInput
              className="block w-12"
              id="reference"
              type="number"
              value={reference.endChapter}
              onChange={(endChapter) =>
                setData((data) => ({
                  ...data,
                  reference: { ...reference, endChapter },
                }))
              }
              autoComplete="off"
              step={1}
              min={1}
            />
          </div>
          <div className="pt-7 mr-2">:</div>
          <div>
            <Label secondary>Verse</Label>
            <NumberInput
              className="block w-12"
              id="reference"
              type="number"
              value={reference.endVerse}
              onChange={(endVerse) =>
                setData((data) => ({
                  ...data,
                  reference: { ...reference, endVerse },
                }))
              }
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
            onChange={(text) => setData((data) => ({ ...data, text }))}
          />
        </div>
        <div className="mt-4">
          <Button type="submit">Save</Button>
        </div>
      </form>
    </Page>
  );
}
