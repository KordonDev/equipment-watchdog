<script lang="ts">
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { getTranslatedGroups } from "../../components/groupsStore";
  import {
    EquipmentType,
    getFreeEquipment,
    type Equipment,
    type EquipmentByType,
  } from "../equipment/equipment.service";
  import type { Equipments, Member } from "./member.service";

  export let member: Member;
  export let submitText: string;
  export let loading: boolean;
  export let onSubmit: (m: Member) => void;

  let allGroups = [];
  getTranslatedGroups.subscribe((groups) => (allGroups = groups));
  let mocked = false;

  onMount(() => {
    member.equipments = addNoneEquipment(member.equipments);
    mocked = true;
  });

  let freeEquipment;
  getFreeEquipment().then((fe) => {
    freeEquipment = allEquipmentPossibilities(member.equipments, fe);
  });

  function handleSubmit() {
    onSubmit({
      ...member,
      equipments: getEquipmentFromFree(member.equipments, freeEquipment),
    });
  }
  const NoEquipment: Equipment = {
    id: 0,
    registrationCode: "Keine",
    type: EquipmentType.TShirt,
  };

  function addNoneEquipment(equipments: Equipments): Equipments {
    return (Object.values(EquipmentType) as EquipmentType[]).reduce(
      (acc: Equipments, eT: EquipmentType) => {
        if (equipments[eT]) {
          acc[eT] = equipments[eT];
        } else {
          acc[eT] = { ...NoEquipment };
        }
        return acc;
      },
      {} as Equipments
    );
  }

  function getEquipmentFromFree(
    selected: Equipments,
    equipments: EquipmentByType
  ): Equipments {
    return (Object.values(EquipmentType) as EquipmentType[]).reduce(
      (acc, eT) => {
        const selectedEquipment = equipments[eT].find(
          (e) => selected[eT].registrationCode === e.registrationCode
        );
        if (
          selected[eT]?.registrationCode !== NoEquipment.registrationCode &&
          selectedEquipment
        ) {
          acc[eT] = selectedEquipment;
        }
        return acc;
      },
      {} as Equipments
    );
  }

  function allEquipmentPossibilities(
    equipments: Equipments,
    equipmentsList: EquipmentByType
  ): EquipmentByType {
    return (Object.values(EquipmentType) as EquipmentType[]).reduce(
      (acc, eT) => {
        acc[eT] = [NoEquipment];
        if (
          equipments[eT] &&
          equipments[eT].registrationCode !== NoEquipment.registrationCode
        ) {
          acc[eT] = acc[eT].concat(equipments[eT]);
        }
        if (equipmentsList[eT]) {
          acc[eT] = acc[eT].concat(...equipmentsList[eT]);
        }
        return acc;
      },
      {} as EquipmentByType
    );
  }
</script>

<form on:submit|preventDefault={handleSubmit} disabled={loading}>
  <Label for="name" class="block mb-2">Name</Label>
  <Input required class="mb-4" id="name" bind:value={member.name} />

  <Label class="mb-4">
    <div class="mb-2">Gruppe</div>
    <Select required items={allGroups} bind:value={member.group} />
  </Label>
  {#if mocked && freeEquipment}
    <div>
      <h2>Ausrüstung</h2>
      <Label class="mb-4">
        <div class="mb-2">Jacke</div>
        <Select
          items={freeEquipment[EquipmentType.Jacket]?.map((e) => ({
            value: e.registrationCode,
            name: e.registrationCode,
          })) || []}
          bind:value={member.equipments[EquipmentType.Jacket].registrationCode}
        />
      </Label>
    </div>
  {/if}
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
