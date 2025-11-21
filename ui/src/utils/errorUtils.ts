export function countErrors<T extends object>(errors: T | undefined): number {
  return Object.values(errors ?? {}).reduce<number>((count, err) => {
    if (Array.isArray(err)) return count + err.length;
    return count;
  }, 0);
}
