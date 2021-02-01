import React from "react";
import {Meta, Story} from "@storybook/react";

import TermInput, {TermInputProps} from ".";

export default {
  title: "TermInput",
  component: TermInput
} as Meta<TermInputProps>;

export const Default: Story<TermInputProps> = (props) => <TermInput {...props} />;
