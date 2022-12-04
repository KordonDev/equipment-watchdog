<script lang="ts">
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
  } from "flowbite-svelte";

  export let users: User[];
  export let onApprove: (userId: number) => void;
  export let onAdmin: (userId: number) => void;
</script>

<Table>
  <TableHead>
    <TableHeadCell>Name</TableHeadCell>
    <TableHeadCell>Approved</TableHeadCell>
    <TableHeadCell>Admin</TableHeadCell>
  </TableHead>
  <TableBody>
    {#each users as user}
      <TableBodyRow>
        <TableBodyCell>{user.name}</TableBodyCell>
        <TableBodyCell>
          <input
            type="checkbox"
            checked={user.isApproved}
            on:change={() => onApprove(user.id)}
          />
        </TableBodyCell>
        <TableBodyCell>
          <input
            type="checkbox"
            checked={user.isAdmin}
            on:change={() => onAdmin(user.id)}
          />
        </TableBodyCell>
      </TableBodyRow>
    {/each}
  </TableBody>
</Table>
