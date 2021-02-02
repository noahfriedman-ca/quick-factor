import React from "react";
import TermInputField from "./Field";

export interface TermInputProps {
  onSubmit?: (values: {[exponent: number]: number}) => void
}

const TermInput: ((props: TermInputProps) => JSX.Element) & {
  Field: typeof TermInputField
} = ({onSubmit}) => (
  <div>
  </div>
);

TermInput.Field = TermInputField;

export default TermInput;
