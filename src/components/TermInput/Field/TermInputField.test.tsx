import React from "react";
import {render, RenderResult} from "@testing-library/react";

import TermInputField from ".";

describe("the TermInput.Field component", () => {
  let r: RenderResult;
  beforeEach(() => {
    r = render(<TermInputField/>);
  });

  test("placeholder", () => {});
});
