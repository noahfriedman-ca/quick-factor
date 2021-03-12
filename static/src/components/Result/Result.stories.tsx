import React from "react";
import {Meta, Story} from "@storybook/react";

import Result, {ResultProps} from ".";

interface Props extends Omit<ResultProps, "factored"> {
  expression?: string
  intercepts?: string[]
}

export default {
  title: "Result",
  component: Result,
  args: {
    result: "full",
    expression: "(x + 5)(x - 4)",
    intercepts: ["-5", "4"],
  },
  argTypes: {
    result: {
      control: {
        type: "select",
        options: [
          "full",
          "quadratic",
          "partial",
          "not",
          "error"
        ]
      },
      table: {
        disable: true
      }
    },
    expression: {
      control: "text"
    },
    intercepts: {
      control: {
        type: "array",
        separator: ", "
      }
    }
  }
} as Meta<Props>;

export const Default: Story<Props> = ({result, expression, intercepts}) => <Result result={result} factored={result === "not" || result === "error" ? undefined : {expression: expression || "", intercepts: intercepts || []}}/>;
Default.argTypes = {
  result: {
    table: {
      disable: false
    }
  }
};

export const Full = Default.bind({});
Full.args = {
  result: "full"
};

export const Quadratic = Default.bind({});
Quadratic.args = {
  result: "quadratic"
};

export const Partial = Default.bind({});
Partial.args = {
  result: "partial"
};


export const Not = Default.bind({});
Not.args = {
  result: "not"
};
Not.argTypes = {
  expression: {
    table: {
      disable: true
    }
  },
  intercepts: {
    table: {
      disable: true
    }
  }
};

export const Error = Default.bind({});
Error.args = {
  result: "error"
};
Error.argTypes = {
  expression: {
    table: {
      disable: true
    }
  },
  intercepts: {
    table: {
      disable: true
    }
  }
};
