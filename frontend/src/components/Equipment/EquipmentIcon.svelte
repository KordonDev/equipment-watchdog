<script lang="ts">
  import type { EquipmentType } from "src/views/equipment/equipment.service";

  export let registrationCode: string | undefined;
  export let equipmentType: EquipmentType | undefined;
</script>

<div>
  <img
    src={`/${equipmentType}.svg`}
    alt=""
    class={`equipment ${registrationCode?.length > 0 ? "equipped" : ""}`}
  />
  {registrationCode ? registrationCode : ""}
</div>

<style>
  img.equipment {
    width: 100px;
    height: 100px;
    padding: 10px;
    margin: 5px;
    border-radius: 5px;
  }
  img.equipment.equipped {
    background: rgba(0, 100, 0, 0.4);
  }
  div {
    color: wheat;
    margin: 10px;
  }
</style>
