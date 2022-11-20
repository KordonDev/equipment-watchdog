<script lang="ts">
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import { groupsStore } from "../../components/groupsStore";
  import type { Member } from "./member.service";

  export let member: Member;
  export let submitText: string;
  export let loading: boolean;
  export let onSubmit: (m: Member) => void;

  let allGroups = [];
  groupsStore.subscribe((value) => (allGroups = value));

  function handleSubmit() {
    onSubmit(member);
  }
</script>

<form on:submit|preventDefault={handleSubmit} disabled={loading}>
  <Label for="name" class="block mb-2">Name</Label>
  <Input required class="mb-4" id="name" bind:value={member.name} />

  <Label class="mb-4">
    <div class="mb-2">Gruppe</div>
    <Select required items={allGroups} bind:value={member.group} />
  </Label>
  <div class="flex flex-row justify-end">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="purple" type="submit">
        {submitText}
      </Button>
    {/if}
  </div>
</form>
