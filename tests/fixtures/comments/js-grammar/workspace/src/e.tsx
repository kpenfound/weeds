/** Props are the props. */
export interface Props {
  /** Label is shown. */
  label: string;
}

/** Button renders a button. */
export const Button = (p: Props) => <button>{p.label}</button>; // trailing
