import React from "react";
import { Card, CardBody, CardText } from "reactstrap";

export class Report extends React.Component {
  render() {
    return (
      <Card className="report">
        <CardBody>
          <CardText>{this.props.report.Description}</CardText>
        </CardBody>
      </Card>
    );
  }
}
