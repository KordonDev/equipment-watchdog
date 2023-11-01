import { fetchApi } from "../apiService";


export function getChangesForEquipment(id: number): Promise<string[]> {
  return fetchApi(`/changes/equipments/${id}`);
}

export function getChangesForOrder(id: number): Promise<string[]> {
  return fetchApi(`/changes/orders/${id}`);
}

export function getChangesForMember(id: number): Promise<string[]> {
  return fetchApi(`/changes/members/${id}`);
}
