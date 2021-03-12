import React from "react";
import {Card} from "react-bootstrap";
import {MathComponent} from "mathjax-react";

export interface ResultProps {
  result: "full" | "quadratic" | "partial" | "not" | "error"
  factored?: {
    expression: string
    intercepts: string[]
  }
}

const Result: (props: ResultProps) => JSX.Element = ({result, factored}) => {
  const [niceResult, setNiceResult] = React.useState<string>();

  React.useEffect(() => {
    switch (result) {
      case "full":
        setNiceResult("The polynomial was fully factored.");
        break;
      case "quadratic":
        setNiceResult("The polynomial was fully factored using the quadratic formula.");
        break;
      case "partial":
        setNiceResult("The polynomial was partially factored.");
        break;
      case "not":
        setNiceResult("The polynomial could not be factored.");
        break;
      case "error":
        setNiceResult("An error occurred and the polynomial could not be factored.");
        break;
    }
  }, [result]);

  return (
    <Card bg={result === "error" || result === "not" ? "primary" : result === "partial" ? "warning" : "success"} text={result === "error" || result === "not" ? "white" : undefined}>
      <Card.Header style={{display: "flex", alignItems: "center"}}>
        <h4 style={{margin: 0, display: "inline"}}>RESULT:</h4>
        <p style={{marginLeft: 8, marginBottom: 0, display: "inline"}}>{niceResult}</p>
      </Card.Header>
      {(result !== "error" && result !== "not") && (
        <Card.Body>
          <b>Expression:</b>
          <span style={{display: "inline-flex"}}><MathComponent tex={factored?.expression} /></span>
          <br />
          <br />
          <b>Intercepts:</b>
          <span style={{display: "inline-flex"}}>
            <MathComponent tex={factored?.intercepts.join(", ")}/>
          </span>
        </Card.Body>
      )}
    </Card>
  );
};

export default Result;
