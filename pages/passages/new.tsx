import { useRouter } from "next/router";
import Page from "../../components/ui/Page";
import PageHeader from "../../components/ui/PageHeader";
import PageTitle from "../../components/ui/PageTitle";
import BackLink from "../../components/ui/BackLink";
import PassageForm from "../../components/PassageForm";

export default function NewPassagePage() {
  const router = useRouter()

  async function onSubmit(data: { text: string, reference: string }) {
    await fetch("/api/passages", {
      method: "POST",
      body: JSON.stringify(data),
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
      <PassageForm onSubmit={onSubmit}/>
    </Page>
  );
}
