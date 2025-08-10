<script lang="ts">
  import {Button, ButtonGroup, Input, Label, Select} from "flowbite-svelte";
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import {
    getTranslatedGroups,
    groupFilter,
  } from "../../components/groupsStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { getMembers, type Member } from "./member.service";
  import MemberCard from "./MemberCard.svelte";

  let allGroups = [];
  getTranslatedGroups.subscribe((groups) => {
    if (groups) {
      allGroups = [{ name: "Alle", value: "all" }, ...groups];
    }
  });

  let groupFilterInternal;
  groupFilter.subscribe((v) => (groupFilterInternal = v));
  function updateGroupFilter(group: string) {
    groupFilterInternal = group;
    groupFilter.set(group);
  }

  function byGroup(m: Member, filter?: string) {
    if (!filter || filter === "all") {
      return true;
    }
    return m.group === filter;
  }

  function byName(m: Member, search: string) {
    if (search === "") {
      return true;
    }
    return m.name.includes(search);
  }

  let search = "";

  let membersPromise = getMembers();
</script>

<Navigation />
<div class="my-2">
  <a href={routes.AddMember.link} use:link>Mitglied hinzufügen</a>
</div>
<Label class="my-4" size="md">
  <div class="mb-2">Gruppe</div>
  <ButtonGroup class="d-flex flex-wrap">
    {#each allGroups as group}
      <Button
        class="m-1"
        on:click={() => updateGroupFilter(group.value)}
        color={groupFilterInternal === group.value ? 'green' : 'purple'}
      >
        {group.name}
      </Button>
    {/each}
  </ButtonGroup>
</Label>
<Label class="block mb-2">
  Suche
  <Input type="search" class="mb-4" bind:value={search} />
</Label>
{#await membersPromise then members}
  <div class="flex flex-wrap">
    {#each members
      .filter((m) => byGroup(m, groupFilterInternal))
      .filter((m) => byName(m, search)) as member}
      <MemberCard {member} columns={2} />
    {/each}
  </div>
{/await}
