import { fetchApi } from "../apiService";

export enum ChangeAction {
	CreateOrder = "create-order",
	DeleteOrder = "delete-order",
	UpdateOrder = "update-order",
	OrderToEquipment = "order-to-equipment",
	CreateMember = "create-member",
	UpdateMember = "update-member",
	DeleteMember = "delete-member",
	CreateEquipment = "create-equipment",
	DeleteEquipment = "delete-equipment"
}

export interface Change {
	id: number;
	createdAt: string;
	memberId: number;
	equipmentId: number;
	orderId: number;
	action: ChangeAction;
	userId: number;
}



export function getChangesForEquipment(id: number): Promise<string[]> {
  return fetchApi(`/changes/equipments/${id}`);
}

export function getChangesForOrder(id: number): Promise<string[]> {
  return fetchApi(`/changes/orders/${id}`);
}

export function getChangesForMember(id: number): Promise<string[]> {
  return fetchApi(`/changes/members/${id}`);
}

export function getAllChanges(): Promise<string[]> {
	return fetchApi(`/changes/`);
}

export function getRecentChanges(): Promise<Change[]> {
	return fetchApi(`/changes/recent`);
}
