export const getDate = (date: string): Date | undefined => {
  if (date === "0001-01-01T00:00:00Z") {
    return undefined;
  }
  return new Date(date);
}

export const formatToDate = (date?: Date | string): string => {
  if (!date) {
    return "-";
  }
	const d = typeof date === 'string' ? new Date(date) : date;
  return d.toLocaleDateString('de-de', { year: "numeric", month: "short", day: "numeric" });
}

