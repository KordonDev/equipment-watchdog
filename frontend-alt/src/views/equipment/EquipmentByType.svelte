<script lang="ts">
  import { routes } from "../../routes";
  import { link, push } from "svelte-spa-router";
  import { afterUpdate, onMount } from "svelte/internal";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import {
    EquipmentType,
    getEquipmentByType,
    translatedEquipmentTypes,
    translateEquipmentType,
    type Equipment
  } from "./equipment.service";
  import {Alert, Card, Label, Spinner, Input, ButtonGroup, Button} from "flowbite-svelte";
  import EquipmentIcon from "../../components/Equipment/EquipmentIcon.svelte";
  import { getMembers, type Member} from "../member/member.service";

  export let params = { type: undefined };
  let currentType;
  let search = "";

  onMount(() => {
    currentType = params.type;
    if (translateEquipmentType(params.type) === params.type) {
      push(`${routes.EquipmentType.link}${EquipmentType.Jacket}`);
    }
  });

  afterUpdate(() => {
    if (currentType !== params.type) {
      currentType = params.type;
      equipmentsPromise = getEquipmentByType(params.type);
    }
  });

  function updateUrl(equipmentType: string) {
    currentType = equipmentType;
    push(`${routes.EquipmentType.link}${equipmentType}`);
  }

  function byCode(e: Equipment, search: string) {
    if (search === "") {
      return true;
    }
    return e.registrationCode.includes(search);
  }

  const findMemberName = (memberId: number, members: Member[]) => {
    return members.find(m => String(m.id) === String(memberId))?.name || "Ausgerüstet";
  }

  let equipmentsPromise = getEquipmentByType(params.type);
  let membersPromise = getMembers();
</script>

<Navigation />

<div class="my-2">
  <a href={`${routes.EquipmentAdd.link}${currentType}`} use:link>Neue Ausrüstung anlegen</a>
</div>

<Label class="mb-4">
  <div class="mb-2">Art</div>
  <ButtonGroup class="d-flex flex-wrap">
    {#each translatedEquipmentTypes as equipment}
      <Button
              class="m-1"
              on:click={() => updateUrl(equipment.value)}
              color={currentType === equipment.value ? 'green' : 'purple'}
      >
        {equipment.name}
      </Button>
    {/each}
  </ButtonGroup>
</Label>
<Label class="block mb-2">
  Suche nach Registrierungsnummer
  <Input type="search" class="mb-4" bind:value={search} />
</Label>
{#await equipmentsPromise}
  <Spinner />
{:then equipments}
  <div class="flex flex-wrap">
    {#each equipments.filter(e => byCode(e, search)) as equipment}
      <div>
        <Card class="m-4" style="padding: 10px">
          <EquipmentIcon
            equipmentType={equipment.type}
            registrationCode={undefined}
          />
          <h3>
            <a href={`${routes.EquipmentDetails.link}${equipment.id}`} use:link>
              {equipment.registrationCode}
            </a>
          </h3>
          <p>Größe: {equipment.size || "-"}</p>
          {#if equipment.memberId && equipment.memberId !== 0}
            {#await membersPromise}
              <Spinner />
            {:then members}
              { findMemberName(equipment.memberId, members) }
            {:catch e}
              Ausgerüstet
            {/await}
          {/if}
        </Card>
      </div>
    {/each}
  </div>
{:catch e}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnte die Ausrüstung nicht geladen werden.
  </Alert>
{/await}
