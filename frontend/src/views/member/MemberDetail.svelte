<script lang="ts">
  import { Button, Modal, Spinner } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { push } from "svelte-spa-router";
  import {
    errorNotification,
    successNotification,
  } from "../../components/Notification/notificationStore";
  import {
    deleteMember,
    getMember,
    updateMember,
    type Member,
  } from "./member.service";
  import MemberForm from "./MemberForm.svelte";
  import MemberCard from "./MemberCard.svelte";
  import Navigation from "../../components/Navigation/Navigation.svelte";

  export let params = { id: undefined };

  let memberPromise = getMember(params.id);
  let deleteModalOpen = false;

  let loading = false;

  function deleteMemberInternal(id: string, name: string) {
    loading = true;
    deleteMember(id)
      .then(() => {
        successNotification(`Mitglied ${name} wurde erfolgreich gelöscht.`);
        loading = false;
        deleteModalOpen = false;
        push(routes.MemberOverview.link);
      })
      .catch(() =>
        errorNotification(`Mitglied ${name} konnte nicht gelöscht werden.`)
      );
  }

  function updateMemberInternal(member: Member) {
    loading = true;
    updateMember(member)
      .then((m) => {
        successNotification(
          `Mitglied ${m.name} wurde erfolgreich gespeichert.`
        );
        // TODO: Fix, should be removed
        //window.location.reload();
      })
      .catch(() =>
        errorNotification(
          `Mitglied ${member.name} konnte nicht gespeichert werden.`
        )
      )
      .finally(() => (loading = false));
  }
</script>

<Navigation />

{#await memberPromise then member}
  <MemberCard {member} columns={6} />

  <MemberForm
    {member}
    onSubmit={updateMemberInternal}
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
