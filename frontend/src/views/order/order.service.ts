import { getDate } from "../../components/timeHelper";
import { fetchApi } from "../apiService";
import type { Equipment, EquipmentType } from "../equipment/equipment.service";
import { useRegistrationCode } from "../registrationCode/registrationCode.service";

export interface Order<T = Date> {
  id: number;
  type: EquipmentType;
  createdAt: T;
  size: string;
  memberId: number;
  fulfilledAt?: T;
}

export function getOrders(): Promise<Order[]> {
  return fetchApi(`/orders/`)
    .then(orders => orders.map(parseOrderDates));
}

export function getFulfilledOrders(): Promise<Order[]> {
  return fetchApi(`/orders/fulfilled`)
    .then(orders => orders.map(parseOrderDates));
}

export function getOrdersForMember(id: string): Promise<Order[]> {
  return fetchApi(`/orders/member/${id}`)
    .then(orders => orders.map(parseOrderDates));
}

export function getOrder(id: string): Promise<Order> {
  return fetchApi(`/orders/${id}`)
    .then(parseOrderDates);
}

export function createOrder(order: Order): Promise<Order> {
  return fetchApi("/orders/", {
    method: "POST",
    body: JSON.stringify({
      ...order,
      id: 0,
    }),
  })
    .then(parseOrderDates);
}

export function updateOrder(order: Order): Promise<Order> {
  return fetchApi(`/orders/${order.id}`, {
    method: "PUT",
    body: JSON.stringify(order),
  })
    .then(parseOrderDates);
}

export function deleteOrder(id: number): Promise<void> {
  return fetchApi(`/orders/${id}`, {
    method: "DELETE",
  });
}

interface Change {
  id: number;
  createdAt: string;
  toMember: number;
  equipment: number;
  order: number;
  action: number;
  byUser: number;
}

export function getAllChanges(): Promise<Change[]> {
  return fetchApi(`/changes/`);
}


export function fulfillOrder(order: Order, registrationCode: string): Promise<Equipment> {
  return fetchApi(`/orders/${registrationCode}/toEquipment`, {
    method: "POST",
    body: JSON.stringify(order)
  })
    .then(e => {
      useRegistrationCode(e.registrationCode);
      return e;
    })
}

function parseOrderDates(order: Order<string>): Order {
  return {
    ...order,
    createdAt: getDate(order.createdAt),
    fulfilledAt: getDate(order.fulfilledAt),
  }

}
