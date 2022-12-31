<script lang="ts">
  import { Button, Modal, Spinner } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { push } from "svelte-spa-router";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { deleteMember, getMember, updateMember } from "./member.service";
  import MemberForm from "./MemberForm.svelte";
  import MemberCard from "./MemberCard.svelte";
  import Navigation from "../../components/Navigation/Navigation.svelte";

  export let params = { id: undefined };

  let memberPromise = getMember(params.id);
  let deleteModalOpen = false;

  let loading = false;

  function deleteMemberInternal(id: string, name: string) {
    deleteMember(id)
      .then(() => {
        createNotification(
          {
            color: "green",
            text: `Mitglied ${name} wurde erfolgreich gelöscht.`,
          },
          5
        );
        loading = false;
        deleteModalOpen = false;
        push(routes.MemberOverview.link);
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Mitglied ${name} konnte nicht gelöscht werden.`,
          },
          20
        );
        loading = false;
      });
  }
</script>

<Navigation />

{#await memberPromise then member}
  <MemberCard {member} columns={6} />

  <MemberForm
    {member}
    onSubmit={updateMember}
    submitText="Speichern"
    loading={false}
    hideEquipment={false}
  />

  <Button color="red" on:click={() => (deleteModalOpen = true)}>
    Mitglied löschen
  </Button>
  <Modal
    title={`Mitglied "${member.name}" löschen?`}
    bind:open={deleteModalOpen}
    autoclose
  >
    <p>
      Wird ein Mitglied gelöscht werden auch auch die zugewiese Ausrüstung der
      Kleiderkammer zugewiesen.
    </p>
    <svelte:fragment slot="footer">
      {#if loading}
        <Spinner />
      {:else}
        <Button
          color="red"
          on:click={() => deleteMemberInternal(member.id, member.name)}
        >
          Löschen
        </Button>
        <Button color="purple">Abbrechen</Button>
      {/if}
    </svelte:fragment>
  </Modal>
{/await}
