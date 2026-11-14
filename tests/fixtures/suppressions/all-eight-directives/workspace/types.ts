export const x: number = 1;
// @ts-expect-error wrong on purpose
const y: string = x;
// @ts-ignore
const z: string = x;
