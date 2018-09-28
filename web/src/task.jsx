import React from "react";
import {
  Card,
  CardBody,
  CardTitle,
  CardImg,
  CardText,
  CardGroup,
  Progress
} from "reactstrap";
import { Report } from "./report";
import { CreateReportModal } from "./createReportModal";

export class Task extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showCreateReportModal: false
    };

    this.toggleShowCreateReportModaltoggle = this.toggleShowCreateReportModal.bind(
      this
    );
  }

  toggleShowCreateReportModal() {
    this.setState({
      showCreateReportModal: !this.state.showCreateReportModal
    });
  }

  render() {
    var reports = [];
    var taskCompletion = 0;
    if (this.props.task.Reports != null) {
      this.props.task.Reports.forEach(function(item) {
        reports.push(<Report key={item.ID} report={item} />);
      });
      taskCompletion =
        (this.props.task.Reports.length /
          (this.props.task.Reports.length + 1)) *
        100;
    }

    return (
      <div className="task" id={this.props.task.ID}>
        <CreateReportModal task={this.props.task} />
        <Card
          body
          inverse
          color="primary"
          style={{ maxWidth: "250px", minWidth: "250px" }}>
          <CardBody>
            <CardTitle>{this.props.task.Name}</CardTitle>
            <CardText>{this.props.task.Description}</CardText>
          </CardBody>
        </Card>
        <CardGroup className="report-container">
          <CardGroup className="report-list">{reports}</CardGroup>
          <Card
            body
            inverse
            color="success"
            style={{ maxWidth: "75px", minWidth: "75px" }}>
            <a
              onClick={this.toggleShowCreateReportModal}
              style={{ margin: "auto" }}>
              <CardImg
                src="https://upload.wikimedia.org/wikipedia/commons/c/c3/Android_Emoji_2795.svg"
                alt="Create New Report"
              />
            </a>
          </Card>
        </CardGroup>
        <div className="progress">
          <Progress value={taskCompletion} />
        </div>
      </div>
    );
  }
}
