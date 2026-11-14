// TODO: remove the global.
export const global = 1;

/** Run starts work. */
export function run(): void {
  // body comment
  const x = 1;
  console.log(x);
}

/** Helper is hidden. */
function helper(): void {}

// const old = compute();
// return old;
export const gone = 2;

export const trailing = 3; // trailing note

/** Detached is separated by a blank line. */

export function detached(): void {}

// FIXME: broken
export function broken(): void {}

export class Config {
  /** Name is exported. */
  name = "";
  /** Secret is private. */
  private secret = "";
  /** Token is protected. */
  protected token = "";
  /** Hash is a hash private. */
  #hash = "";
}

// plain comment before an exported function
export function plain(): void {}

class Hidden {
  /** Member of an unexported class. */
  a = 1;
}
