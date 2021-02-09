import React from "react";
import {render, RenderResult, fireEvent} from "@testing-library/react";

import TermInput from ".";

describe("the TermInput component", () => {
  let r: RenderResult;
  beforeEach(() => {
    r = render(<TermInput/>);
  });

  it("should not display any TermInput.Field components until a number of terms is entered", () => {
    const i = r.getByRole("spinbutton"); // Throws error unless exactly one matching element is found

    fireEvent.change(i, {target: {value: "10"}});
    fireEvent.click(r.getByText(/go/i));

    expect(r.getAllByRole("spinbutton")).toHaveLength(12);
  });
  describe("should not acknowledge input and display an error when a value other than a positive integer that is >= 3 is given", () => {
    test.each(["3.14", "-2", "0", "1"])("value %s", (v: string) => {
      const i = r.getByRole("spinbutton");

      fireEvent.change(i, {target: {value: "4"}});
      fireEvent.click(r.getByText(/go/i));

      expect(r.queryByText(/^error/i)).toBeNull();

      fireEvent.change(i, {target: {value: v}});
      fireEvent.click(r.getByText(/go/i));

      r.getByText(/^error/i);
      expect(r.getAllByRole("spinbutton")).toHaveLength(6);
    });
  });
  it("should respond with an object mapping values to their corresponding exponents when the form is submitted", () => {
    const mockEvent = jest.fn();
    r.rerender(<TermInput onSubmit={mockEvent} />);

    fireEvent.change(r.getByRole("spinbutton"), {target: {value: "3"}});
    fireEvent.click(r.getByText(/go/i));

    r.getAllByRole("spinbutton").forEach((v, i) => {
      if (i === 0) {
        return;
      }

      fireEvent.change(v, {target: {value: `${i}`}});
    });
    fireEvent.submit(r.getByRole("form"));

    expect(mockEvent).toHaveBeenCalledWith([4, 3, 2, 1]);
  });
  it("should use the value '0' when an input field is empty", () => {
    const mockEvent = jest.fn();
    r.rerender(<TermInput onSubmit={mockEvent} />);

    fireEvent.change(r.getByRole("spinbutton"), {target: {value: "2"}});
    fireEvent.click(r.getByText(/go/i));

    fireEvent.submit(r.getByRole("form"));

    expect(mockEvent).toHaveBeenCalledWith([0, 0, 0]);
  });
});
