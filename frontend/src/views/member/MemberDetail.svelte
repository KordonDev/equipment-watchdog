<script lang="ts">
  import { createNotification } from "../../components/notificationStore";
  import { getMember, updateMember } from "./member.service";
  import MemberForm from "./MemberForm.svelte";
  export let params = { id: undefined };

  let memberPromise = getMember(params.id);
</script>

{#await memberPromise then member}
  <h1>{member.name}</h1>
  <MemberForm
    {member}
    onSubmit={updateMember}
    submitText="Speichern"
    loading={false}
  />
{/await}
