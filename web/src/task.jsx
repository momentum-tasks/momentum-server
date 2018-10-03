import React from "react";
import axios from "axios";
import {
  Button,
  Card,
  CardBody,
  CardTitle,
  CardImg,
  CardText,
  CardGroup,
  Form,
  FormGroup,
  Input,
  Label,
  Modal,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Progress
} from "reactstrap";
import { Report } from "./report";

export class Task extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      task: this.props.task,
      showCreateReportModal: false,
      description: ""
    };

    this.toggleShowCreateReportModal = this.toggleShowCreateReportModal.bind(
      this
    );
  }

  toggleShowCreateReportModal() {
    this.setState({
      showCreateReportModal: !this.state.showCreateReportModal,
      sequence: 0,
      description: ""
    });
  }

  validateCreateReportModal() {
    return this.state.sequence > 0 && this.state.description.length > 0;
  }

  handleChange = event => {
    this.setState({
      [event.target.id]: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();
    axios
      .post(
        "http://localhost:3000/tasks/" + this.state.task.ID + "/reports",
        {
          sequence: parseInt(this.state.sequence),
          description: this.state.description
        },
        {
          headers: {
            "X-Session-Token": this.props.token,
            "Content-Type": "application/json"
          }
        }
      )
      .then(response => {
        axios
          .get("http://localhost:3000/tasks/" + this.state.task.ID, {
            headers: { "X-Session-Token": this.props.token }
          })
          .then(response => {
            this.setState({ task: response.data });
            this.toggleShowCreateReportModal();
          });
      })
      .catch(error => {
        console.log(error.response);
      });
  };

  render() {
    var reports = [];
    var taskCompletion = 0;
    if (this.state.task.Reports != null) {
      this.state.task.Reports.forEach(function(item) {
        reports.push(<Report key={item.ID} report={item} />);
      });
      taskCompletion =
        (this.state.task.Reports.length /
          (this.state.task.Reports.length + 1)) *
        100;
    }

    return (
      <div>
        <Modal
          isOpen={this.state.showCreateReportModal}
          toggle={this.toggleShowCreateReportModal}
          className={this.props.className}>
          <ModalHeader toggle={this.toggle}>Modal title</ModalHeader>
          <ModalBody>
            <Form onSubmit={this.handleSubmit}>
              <FormGroup>
                <Label>Sequence</Label>
                <Input
                  autoFocus
                  type="number"
                  id="sequence"
                  onChange={this.handleChange}
                />
              </FormGroup>
              <FormGroup>
                <Label>Description</Label>
                <Input
                  type="textarea"
                  id="description"
                  onChange={this.handleChange}
                />
              </FormGroup>
              <Button
                color="primary"
                disabled={!this.validateCreateReportModal()}>
                Submit
              </Button>
            </Form>
          </ModalBody>
        </Modal>
        <div className="task" id={this.state.task.ID}>
          <Card
            body
            inverse
            color="primary"
            style={{ maxWidth: "250px", minWidth: "250px" }}>
            <CardBody>
              <CardTitle>{this.state.task.Name}</CardTitle>
              <CardText>{this.state.task.Description}</CardText>
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
        </div>
        <Progress color="success" value={taskCompletion} />
      </div>
    );
  }
}
