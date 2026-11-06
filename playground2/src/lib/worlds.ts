/** Clamp a selected world index into [0, count); 0 when there are no worlds. */
export function clampWorldIndex(index: number, count: number): number {
  if (count <= 0 || !Number.isFinite(index)) return 0;
  return Math.min(Math.max(Math.trunc(index), 0), count - 1);
}

export function worldLabel(index: number, kind: string, relationships: string, name?: string): string {
  const n = relationships.split("\n").filter((l) => l.trim() !== "").length;
  return `${name ?? `World ${index + 1}`} (${kind}, ${n} ${n === 1 ? "relationship" : "relationships"})`;
}
