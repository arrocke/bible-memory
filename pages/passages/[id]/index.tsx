import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/router";
import { format, parse } from 'date-fns'
import Page from "../../../components/ui/Page";
import PageHeader from "../../../components/ui/PageHeader";
import PageTitle from "../../../components/ui/PageTitle";
import BackLink from "../../../components/ui/BackLink";
import Label from "../../../components/ui/Label";
import Button from "../../../components/ui/Button";
import MultilineTextInput from "../../../components/ui/MultilineTextInput";
import NumberInput from "../../../components/ui/NumberInput";
import TextInput from "../../../components/ui/TextInput";

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
  level: number
  reviewDate?: string
}

const REFERENCE_MATCH = /^((?:(?:1|2) )?\w+) (\d+):(\d+)-(\d+):(\d+)$/;
function parseReference(reference: string): Reference {
  const match = REFERENCE_MATCH.exec(reference);
  if (match) {
    return {
      book: match[1],
      startChapter: parseInt(match[2]),
      startVerse: parseInt(match[3]),
      endChapter: parseInt(match[4]),
      endVerse: parseInt(match[5]),
    };
  } else {
    return {};
  }
}

export default function EditPassagePage() {
  const router = useRouter();
  const [passage, setPassage] = useState<PassageFormData>();
  const id = router.query.id as string;

  const referenceString = passage
    ? `${passage.reference.book} ${passage.reference.startChapter}:${passage.reference.startVerse}-${passage.reference.endChapter}:${passage.reference.endVerse}`
    : "";

  useEffect(() => {
    if (typeof id === "string") {
      loadPassage(id);
    }
  }, [id]);

  async function loadPassage(id: string) {
    const request = await fetch(`/api/passages/${id}`);
    const passage = await request.json();
    setPassage({
      reference: parseReference(passage.reference),
      text: passage.text,
      level: passage.level,
      reviewDate: passage.reviewDate ? format(new Date(passage.reviewDate), 'MM/dd/yyyy') : undefined
    });
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (passage) {
      await fetch(`/api/passages/${id}`, {
        method: "PATCH",
        body: JSON.stringify({
          reference: referenceString,
          text: passage.text,
          level: passage.level,
          reviewDate: passage.reviewDate ? parse(passage.reviewDate, 'MM/dd/yyyy', new Date()) : undefined
        }),
        headers: {
          "content-type": "application/js",
        },
      });
      router.push("/passages");
    }
  }

  async function onDelete() {
    await fetch(`/api/passages/${id}`, {
      method: "DELETE",
    });
    router.push("/passages");
  }

  return (
    <Page>
      <PageHeader>
        <BackLink href="/passages">Back to Passages</BackLink>
        <PageTitle>Edit {referenceString}</PageTitle>
      </PageHeader>
      {passage ? (
        <form onSubmit={onSubmit}>
          <fieldset className="flex">
            <Label as="legend">REFERENCE</Label>
            <div className="mr-6">
              <Label secondary>BOOK</Label>
              <TextInput
                className="block w-32"
                id="reference"
                value={passage.reference.book}
                onChange={(book) =>
                  setPassage({
                    ...passage,
                    reference: { ...passage.reference, book },
                  })
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
                value={passage.reference.startChapter}
                onChange={(startChapter) =>
                  setPassage({
                    ...passage,
                    reference: { ...passage.reference, startChapter },
                  })
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
                value={passage.reference.startVerse}
                onChange={(startVerse) =>
                  setPassage({
                    ...passage,
                    reference: { ...passage.reference, startVerse },
                  })
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
                value={passage.reference.endChapter}
                onChange={(endChapter) =>
                  setPassage({
                    ...passage,
                    reference: { ...passage, endChapter },
                  })
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
                value={passage.reference.endVerse}
                onChange={(endVerse) =>
                  setPassage({
                    ...passage,
                    reference: { ...passage.reference, endVerse },
                  })
                }
                autoComplete="off"
                step={1}
                min={1}
              />
            </div>
          </fieldset>
          <div className="mt-2">
            <Label htmlFor="level">LEVEL</Label>
            <NumberInput
              className="block w-12"
              id="level"
              value={passage.level}
              onChange={(level) => setPassage({ ...passage, level: level ?? 0 })}
            />
          </div>
          <div className="mt-2">
            <Label htmlFor="date">REVIEW DATE</Label>
            <TextInput
              className="block w-32"
              id="level"
              value={passage.reviewDate}
              onChange={(reviewDate) => setPassage({ ...passage, reviewDate })}
              pattern="\d{2}\/\d{2}\/\d{4}"
              placeholder="mm/dd/yyyy"
            />
          </div>
          <div className="mt-2">
            <Label htmlFor="text">TEXT</Label>
            <MultilineTextInput
              className="block w-full h-96 resize-none"
              id="text"
              value={passage.text}
              onChange={(text) => setPassage({ ...passage, text })}
            />
          </div>
          <div className="mt-4">
            <Button type="submit">Save</Button>
            <Button className="ml-4" destructive onClick={onDelete}>
              Delete
            </Button>
          </div>
        </form>
      ) : null}
    </Page>
  );
}
