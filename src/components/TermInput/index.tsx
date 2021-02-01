import React from "react";
import TermInputField from "./Field";

export interface TermInputProps {

}

const TermInput: ((props: TermInputProps) => JSX.Element) & {
  Field: typeof TermInputField
} = () => (
  <div>
  </div>
);

TermInput.Field = TermInputField;

export default TermInput;
