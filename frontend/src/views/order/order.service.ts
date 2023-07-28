import { fetchApi } from "../apiService";
import type { EquipmentType } from "../equipment//equipment.service";

export interface Order {
  id: number;
  type: EquipmentType;
  createdAt: string;
  size: string;
  memberId: number;
  fulfilledAt?: string;
}

export function getOrders(): Promise<Order[]> {
  return fetchApi(`/orders/`);
}

export function getFulfilledOrders(): Promise<Order[]> {
  return fetchApi(`/orders/fulfilled`);
}

export function getOrder(id: string): Promise<Order> {
  return fetchApi(`/orders/${id}`);
}

export function createOrder(order: Order): Promise<Order> {
  return fetchApi("/orders/", {
    method: "POST",
    body: JSON.stringify({
      ...order,
      id: 0,
    }),
  });
}

export function updateOrder(order: Order): Promise<Order> {
  return fetchApi(`/orders/${order.id}`, {
    method: "PUT",
    body: JSON.stringify(order),
  });
}

export function deleteOrder(id: number): Promise<void> {
  return fetchApi(`/orders/${id}`, {
    method: "DELETE",
  });
}
