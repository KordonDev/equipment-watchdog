<script lang="ts">
  import { routes } from "../../routes";
  import { push } from "svelte-spa-router";
  import { createNotification } from "../../components/notificationStore";
  import type { Member } from "./member.service";
  import { createMember } from "./member.service";
  import MemberForm from "./MemberForm.svelte";

  let member: Member = {
    id: "0",
    name: "",
    group: "",
  };

  let loading = false;

  function createMemberInternal(m: Member) {
    createMember(m)
      .then((newMember) => {
        createNotification(
          {
            color: "green",
            text: `Mitglied ${m.name} wurde erfolgreich angelegt.`,
          },
          5
        );
        loading = false;
        push(`${routes.MemberDetail.link}${newMember.id}`);
      })
      .catch(() => {
        createNotification({
          color: "red",
          text: `Mitglied ${m.name} konnte nicht angelegt werden.`,
        });
        loading = false;
      });
  }
</script>

<h1>Mitglied hinzufügen</h1>
<MemberForm
  {member}
  onSubmit={createMemberInternal}
  submitText="Anlegen"
  {loading}
/>
