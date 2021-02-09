import React from "react";
import {Meta, Story} from "@storybook/react";

import TermInput, {TermInputProps} from ".";
import {action} from "@storybook/addon-actions";

export default {
  title: "TermInput",
  component: TermInput,
  args: {
    onSubmit: action("form submitted")
  },
  argTypes: {
    onSubmit: {
      table: {
        disable: true
      }
    }
  },
  parameters: {
    controls: {
      hideNoControlsWarning: true
    }
  }
} as Meta<TermInputProps>;

export const Default: Story<TermInputProps> = (props) => <TermInput {...props} />;
