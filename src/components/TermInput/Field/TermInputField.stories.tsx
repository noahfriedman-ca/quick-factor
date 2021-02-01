import React from "react";
import {Meta, Story} from "@storybook/react";

import TermInputField, {TermInputFieldProps} from ".";

export default {
  title: "TermInput/Field",
  component: TermInputField,
  argTypes: {
    exponent: {
      control: {
        type: "number",
        min: "0"
      }
    }
  }
} as Meta<TermInputFieldProps>;

export const Default: Story<TermInputFieldProps> = (props) => <TermInputField exponent={props.exponent ? props.exponent : 0} />;
