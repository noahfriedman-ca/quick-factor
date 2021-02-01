import React from "react";
import {render, RenderResult} from "@testing-library/react";

import TermInput from ".";

describe("the TermInput component", () => {
  let r: RenderResult;
  beforeEach(() => {
    r = render(<TermInput/>);
  });

  test("placeholder", () => {});
});
