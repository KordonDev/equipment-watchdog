<script lang="ts">
  import { Button, Modal, Spinner, Alert, Heading } from "flowbite-svelte";
  import { routes } from "../../routes";
  import { push, link } from "svelte-spa-router";
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
  import { writable } from "svelte/store";
  import OrderCard from "../order/OrderCard.svelte";
  import { getOrdersForMember } from "../order/order.service"

  export let params = { id: undefined };
  const member = writable<Member | undefined>();

  const memberPromise = getMember(params.id)
    .then((m) => member.set(m))

  const ordersPromise = getOrdersForMember(params.id);
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

{#await memberPromise}
  <Spinner />
{:then} 
  
  <div class="mb-4">
    <MemberCard member={$member} columns={6} />
    <MemberForm
      memberStore={member}
      onSubmit={updateMemberInternal}
      submitText="Speichern"
      loading={false}
      hideEquipment={false}
    />
  </div>

  <div class="mb-4">
    <Heading tag="h3">Bestellungen</Heading>
    <a href={`${routes.AddOrder.link}${params.id}`} use:link>Ausrüstung bestellen</a>
    {#await ordersPromise }
      <Spinner />
    {:then orders}
      <div class="flex flex-wrap">
        {#each orders as order}
          <OrderCard order={order} />
        {/each}
      </div>
    {/await}
  </div>

  <Button color="red" on:click={() => (deleteModalOpen = true)}>
    Mitglied löschen
  </Button>
  <Modal
    title={`Mitglied "${$member.name}" löschen?`}
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
          on:click={() => deleteMemberInternal($member.id, $member.name)}
        >
          Löschen
        </Button>
        <Button color="purple">Abbrechen</Button>
      {/if}
    </svelte:fragment>
  </Modal>
{:catch}
  <Alert class="mb-4" color="red">Mitglied konnte nicht geladen werden.</Alert>
{/await}
