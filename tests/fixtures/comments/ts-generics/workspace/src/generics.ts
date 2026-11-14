/** Cast casts a value. */
export function cast<T>(x: unknown): T {
  return <T>x;
}

/** Identity returns its argument. */
export const identity = <T>(x: T): T => x;

/** Wrap wraps a value. */
export const wrap = <T,>(x: T) => [x];

/** Pick picks a value. */
export function pick<K extends string>(key: K): K {
  const k = <K>key;
  return k;
}
