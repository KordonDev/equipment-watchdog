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

export function getRecentChanges(): Promise<Change[]> {
	return fetchApi(`/changes/recent`);
}
