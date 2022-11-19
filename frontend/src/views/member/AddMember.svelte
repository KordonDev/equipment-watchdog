<script lang="ts">
  import { Alert } from "flowbite-svelte";
  import type { Member } from "./member.service";
  import { createMember } from "./member.service";
  import MemberForm from "./MemberForm.svelte";

  let member: Member = {
    id: "0",
    name: "",
    group: "",
  };

  let errorAlert = false;
  let loading = true;

  function createMemberInternal(m: Member) {
    createMember(m)
      .then((succ) => {
        loading = false;
      })
      .catch((err) => {
        errorAlert = true;
        loading = false;
      });
  }

  function closeAlert() {
    errorAlert = false;
  }
</script>

<h1>Mitglied hinzufügen</h1>
<MemberForm {member} onSubmit={createMember} submitText="Anlegen" />
{#if errorAlert}
  <Alert color="red" dismissable on:close={closeAlert}>Close me</Alert>
{/if}
