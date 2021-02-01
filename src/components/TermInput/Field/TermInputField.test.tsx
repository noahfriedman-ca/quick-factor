import React from "react";
import {render, RenderResult} from "@testing-library/react";

import TermInputField from ".";

jest.mock("mathjax-react", () => ({
  MathComponent: (props: {tex?: string}) => <>{props.tex}</>
}));

describe("the TermInput.Field component", () => {
  let r: RenderResult;
  beforeEach(() => {
    r = render(<TermInputField exponent={0} />);
  });

  it("should display the correct exponent based on the 'exponent' prop", () => {
    r.rerender(<TermInputField exponent={2} />);

    r.getByText("x^{2}");
  });
  it("should log a warning if the value in the exponent prop is not an integer, and then round", () => {
    console.warn = jest.fn();

    r.rerender(<TermInputField exponent={3.1415} />);

    r.getByText("x^{3}");
    expect(console.warn).toHaveBeenCalledWith("exponent '3.14' was rounded to '3'");

    (console.warn as jest.Mock).mockRestore();
  });
  it("should log a warning if the value in the exponent prop is not a positive number, and then get the absolute value", () => {
    console.warn = jest.fn();

    r.rerender(<TermInputField exponent={-3} />);

    r.getByText("x^{3}");
    expect(console.warn).toHaveBeenCalledWith("exponent '-3' was flipped to '3'");
  });
});
