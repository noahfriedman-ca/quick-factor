import React from "react";
import {Card, Form} from "react-bootstrap";
import {MathComponent} from "mathjax-react";

import "./TermInputField.css";

export interface TermInputFieldProps {
  exponent: number
  id?: string
}

const TermInputField: (props: TermInputFieldProps) => JSX.Element = ({exponent, id}) => {
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
    <Form.Group style={{margin: 5}}>
      <Card bg="warning" text="dark" style={{flexDirection: "row", padding: 5, alignItems: "center"}}>
        <Form.Control id={id} className="term-input-field" type="number"/>
        {checkedExponent !== 0 && (
          <Form.Label style={{display: "inline-flex", marginLeft: 5, fontSize: "1.2em"}}>
            <MathComponent tex={String.raw`x${checkedExponent !== 1 ? `^{${checkedExponent}}` : ""}`}/>
          </Form.Label>
        )}
      </Card>
    </Form.Group>
  );
};

export default TermInputField;
