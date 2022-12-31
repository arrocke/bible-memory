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

export default function NewPassagePage() {
  const router = useRouter()
  const [reference, setReference] = useState("");
  const [text, setText] = useState("");

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    await fetch("/api/passages", {
      method: "POST",
      body: JSON.stringify({ reference, text }),
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
        <div>
          <Label htmlFor="reference">REFERENCE</Label>
          <TextInput 
            className="block"
            id="reference"
            value={reference}
            onChange={setReference}
            autoComplete="off"
          />
        </div>
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
