export type Notification = {
	id: string;
	message: string;
	type: 'error' | 'info' | 'success';
};

export const notifications = $state<Notification[]>([]);

export function showError(message: string, duration = 3000) {
	showNotification(message, duration, 'error');
}

export function showInfo(message: string, duration = 3000) {
	showNotification(message, duration, 'info');
}
export function showSuccess(message: string, duration = 2000) {
	showNotification(message, duration, 'success');
}

function showNotification(message: string, duration: number, type: 'error' | 'info' | 'success') {
	const id = Math.random().toString(36);
	notifications.push({
		id: id,
		message: message,
		type: type,
	})
	setTimeout(() => {
		const index = notifications.findIndex(n => n.id === id);
		notifications.splice(index, 1)
	}, duration);
}

