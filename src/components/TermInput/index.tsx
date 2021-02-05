import React from "react";
import {Alert, Button, Card, Col, Form, Row} from "react-bootstrap";
import TermInputField from "./Field";
import {MathComponent} from "mathjax-react";

export interface TermInputProps {
  onSubmit?: (values: {[exponent: number]: number}) => void
}

const TermInput: ((props: TermInputProps) => JSX.Element) & {
  Field: typeof TermInputField
} = ({onSubmit}) => {
  const [degree, setDegree] = React.useState(0);
  const [fields, setFields] = React.useState<JSX.Element[]>([]);
  const [error, setError] = React.useState<string>();

  function updateFields() {
    if (degree !== Math.round(degree) || degree < 3) {
      setError("invalid input - must be an integer larger than or equal to 3");
    } else {
      const fields: JSX.Element[] = [];
      for (let i = degree; i >= 0; i--) {
        const inputField = <TermInputField exponent={i} key={`exp${i}`} />;
        if (i !== 0) {
          fields.push(
            <div style={{display: "inline-flex", alignItems: "center"}}>
              {inputField}{" + "}
            </div>
          );
        } else {
          fields.push(inputField);
        }
      }
      setFields(fields);
      setError(undefined);
    }
  }

  return (
    <Card bg="primary" text="light">
      <Card.Header style={{display: "inline-flex", alignItems: "center"}}>
        <label htmlFor="degree" style={{margin: 0}}>
          <Card.Text>What is the degree of the polynomial?</Card.Text>
        </label>
        <input id="degree" type="number" min={3} style={{margin: "0 10px", width: 50}} onChange={v => {
          const n = v.target.valueAsNumber;
          if (isNaN(n)) {
            console.error("Nice try...");
          } else {
            setDegree(n);
          }
        }} />
        <Button onClick={updateFields}>Go!</Button>
        <Alert show={error !== undefined} onClose={() => setError(undefined)} variant="warning" style={{marginLeft: "max(10px, auto)", marginBottom: 0}} className="text-primary" dismissible>
          ERROR: {error}
        </Alert>
      </Card.Header>
      {fields.length !== 0 && (
        <Card.Body>
          <Form inline>
            {fields}
          </Form>
          <hr />
          <Row style={{alignItems: "center"}}>
            <Col md={6}>
              <Button className="btn-block">Factor!</Button>
            </Col>
            <Col md={6}>
              <Alert style={{margin: 0, display: "inline-flex", width: "100%"}} variant="secondary">
                Leave blank for "<MathComponent tex={"0"}/>" value.
              </Alert>
            </Col>
          </Row>
        </Card.Body>
      )}
    </Card>
  );
};

TermInput.Field = TermInputField;

export default TermInput;
