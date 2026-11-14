/*
 * Copyright 2024 Example Corp.
 * Licensed under the MIT licence.
 */

/// <reference types="node" />
// @ts-check

/** Greeter greets people. */
export class Greeter {
  /** Name is who to greet. */
  name = "";
  /** Greet greets. */
  greet(): string {
    return "hello " + this.name;
  }
}

/** Speaker says things. */
export interface Speaker {
  /** Speak says it. */
  speak(): string;
}

/** Mode selects how to greet. */
export enum Mode {
  /** Loud shouts. */
  Loud,
  /** Quiet whispers. */
  Quiet = 1,
}

/** Greet returns a greeting. */
export function greet(name: string): string {
  return `hello ${name}`;
}

/** Version is the library version. */
export const version = "1";

/** Main is the default export. */
export default function main(): void {}

/** Hello re-exports greet. */
export { greet as hello };

/** Legacy is a CommonJS export. */
module.exports.legacy = 1;

/** Keep keeps things. */
// @ts-check
export function keep(): void {}

function helper(): void {
  // @ts-check
  return;
}

//# sourceMappingURL=api.js.map
