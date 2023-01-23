<script lang="ts">
  import { Label, Select } from "flowbite-svelte";
  import {
    NoEquipment,
    translatedEquipmentTypes,
    type Equipment,
    type EquipmentListByType,
    type EquipmentType,
  } from "../../views/equipment/equipment.service";
  import type { Member } from "../..//views/member/member.service";
  import { writable, type Writable } from "svelte/store";

  export let freeEquipments: EquipmentListByType;
  export let memberStore: Writable<Member>;
  export let equipmmentType: EquipmentType;

  const equipmentCode = writable(
    $memberStore.equipments[equipmmentType]?.registrationCode ||
      NoEquipment.registrationCode
  );

  function getItems(freeEquipment: Equipment[]) {
    return (
      freeEquipment?.map((e) => ({
        value: e.registrationCode,
        name: `${e.registrationCode} ${e.size ? "(" + e.size + ")" : ""}`,
      })) || []
    );
  }

  function updateMember() {
    memberStore.update((m) => {
      const selectedEquipment = freeEquipments[equipmmentType]?.find(
        (e) =>
          $equipmentCode === e.registrationCode &&
          $equipmentCode !== NoEquipment.registrationCode
      );
      return {
        ...m,
        equipments: {
          ...m.equipments,
          [equipmmentType]: selectedEquipment,
        },
      };
    });
  }

  function translateType(type: EquipmentType) {
    return translatedEquipmentTypes.find((e) => e.value === type)?.name || type;
  }
</script>

<Label class="mb-4">
  <div class="mb-2">{translateType(equipmmentType)}</div>
  <Select
    items={getItems(freeEquipments[equipmmentType])}
    bind:value={$equipmentCode}
    on:change={updateMember}
  />
</Label>
