import React from "react";
import {Form} from "react-bootstrap";
import {MathComponent} from "mathjax-react";

export interface TermInputFieldProps {
  exponent: number
}

const TermInputField: (props: TermInputFieldProps) => JSX.Element = ({exponent}) => {
  const [checkedExponent, setCheckedExponent] = React.useState<number>();

  React.useEffect(() => {
    let check = exponent;

    const round = Math.round(check);
    if (check !== round) {
      check = round;
      console.warn(`exponent '${parseFloat(exponent.toFixed(2))}' was rounded to '${round}'`);
    }

    const abs = Math.abs(check);
    if (check !== abs) {
      check = abs;
      console.warn(`exponent '${parseFloat(exponent.toFixed(2))}' was flipped to '${abs}'`);
    }

    setCheckedExponent(check);
  }, [exponent]);

  return (
    <Form.Group>
      <Form.Control style={{width: "auto", display: "inline-flex"}} type="number"/>
      {checkedExponent !== 0 && (
        <Form.Label style={{display: "inline-flex", margin: "0 15px 0 5px", fontSize: "1.2em"}}>
          <MathComponent style={{margin: 0}} tex={String.raw`x${checkedExponent !== 1 ? `^{${checkedExponent}}` : ""}`}/>
        </Form.Label>
      )}
    </Form.Group>
  );
};

export default TermInputField;
