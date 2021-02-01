import React from "react";
import {Meta, Story} from "@storybook/react";

import TermInputField, {TermInputFieldProps} from ".";

export default {
  title: "TermInput/Field",
  component: TermInputField
} as Meta<TermInputFieldProps>;

export const Default: Story<TermInputFieldProps> = (props) => <TermInputField {...props} />;
