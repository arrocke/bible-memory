import Link from "next/link";
import { useRouter } from "next/router";
import { FormEvent, useState } from "react";

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
    <div>
      <h1>New Passage</h1>
      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="reference">Reference</label>
          <input
            type="text"
            id="reference"
            required
            value={reference}
            onChange={(e) => setReference(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="text">Text</label>
          <textarea
            id="text"
            required
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
        </div>
        <div>
          <Link href="/passages">Back</Link>
          <button>Save</button>
        </div>
      </form>
    </div>
  );
}
