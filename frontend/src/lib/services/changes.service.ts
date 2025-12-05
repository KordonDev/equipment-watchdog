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
	DeleteEquipment = "delete-equipment",
	UpdateEquipmentOnMember = "update-equipment-on-member"
}

export interface Change {
	id: number;
	createdAt: string;
	memberId: number;
	equipmentId: number;
	oldEquipmentId?: number;
	orderId: number;
	action: ChangeAction;
	userId: number;
}

export function getRecentChanges(): Promise<Change[]> {
	return fetchApi(`/changes/recent`);
}
