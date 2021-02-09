import React from "react";
import {Alert, Button, Card, Col, Form, Row} from "react-bootstrap";
import TermInputField from "./Field";
import {MathComponent} from "mathjax-react";

export interface TermInputProps {
  onSubmit?: (values: number[]) => void
}

const TermInput: ((props: TermInputProps) => JSX.Element) & {
  Field: typeof TermInputField
} = ({onSubmit}) => {
  const [degree, setDegree] = React.useState(0);
  const [fields, setFields] = React.useState<JSX.Element[]>([]);
  const [error, setError] = React.useState<string>();

  function updateFields() {
    if (degree !== Math.round(degree) || degree < 2) {
      setError("invalid input - must be an integer larger than or equal to 2");
    } else {
      const fields: JSX.Element[] = [];
      for (let i = degree; i >= 0; i--) {
        const exp = `exp${i}`;
          fields.push(
            <div style={{display: "inline-flex", alignItems: "center"}} key={exp}>
              <TermInputField exponent={i} id={exp} />{i !== 0 ? " + " : ""}
            </div>
          );
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
          if (isNaN(n) && v.target.value !== "") {
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
          <Form aria-label="form" onSubmit={e => {
            e.preventDefault();
            if (onSubmit) {
              const r: number[] = [];

              for (let v = e.currentTarget.firstElementChild; v !== null; v = v.nextElementSibling) {
                let currentElement = v;
                if (currentElement.tagName === "DIV") {
                  while (currentElement.tagName !== "INPUT") {
                    currentElement = (currentElement.firstElementChild as Element);
                  }
                } else if (currentElement.tagName === "HR") {
                  break;
                }

                let checkValue = (currentElement as HTMLInputElement).valueAsNumber;
                if ((currentElement as HTMLInputElement).value === "") {
                  checkValue = 0;
                }
                if (isNaN(checkValue)) {
                  throw new Error("Nice try...");
                } else {
                  r[parseInt(currentElement.id.replace("exp", ""))] = checkValue;
                }
              }

              onSubmit(r);
            }
          }}>
            {fields}
            <hr />
            <Row style={{alignItems: "center"}}>
              <Col xs={12} sm={4} md={6}>
                <Button className="btn-block" type="submit">Factor!</Button>
              </Col>
              <Col xs={12} sm={8} md={6}>
                <Alert className="mt-3 mt-sm-0" style={{margin: 0, display: "inline-flex", width: "100%"}} variant="secondary">
                  Leave blank for "<MathComponent tex={"0"}/>" value.
                </Alert>
              </Col>
            </Row>
          </Form>
        </Card.Body>
      )}
    </Card>
  );
};

TermInput.Field = TermInputField;

export default TermInput;
