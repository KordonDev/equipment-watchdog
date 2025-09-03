export type Notification = {
	id: string;
	message: string;
	type: 'error' | 'info' | 'success';
};

export const notifications = $state<Notification[]>([]);

export function showError(message: string, duration = 2000) {
	const id = Math.random().toString(36);
	console.log(`${id}: ${message}`);
	notifications.push({
		id: id,
		message: message,
		type: 'error',
	})
	setTimeout(() => {
		const index = notifications.findIndex(n => n.id === id);
		notifications.splice(index, 1)
	}, duration);
}

