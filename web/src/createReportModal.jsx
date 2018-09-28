import React from "react";
import { Button, Form, FormGroup, Label, Input } from "reactstrap";

export class CreateReportModal extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      sequence: 0,
      description: ""
    };
  }

  validateForm() {
    return this.state.sequence > 0 && this.state.description.length > 0;
  }

  handleChange = event => {
    console.log(event.target);
    this.setState({
      [event.target.id]: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();
    axios
      .post("http://localhost:3000/" + this.props.task.ID, null, {
        auth: { username: this.state.username, password: this.state.password }
      })
      .then(response => {
        new Cookies().set("X-Session-Token", response.data);
        window.location.reload();
      });
  };

  render() {
    return (
      <div>
        <Modal
          isOpen={this.state.modal}
          toggle={this.toggle}
          className={this.props.className}>
          <ModalHeader toggle={this.toggle}>Modal title</ModalHeader>
          <ModalBody>asdf</ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={this.toggle}>
              Do Something
            </Button>{" "}
            <Button color="secondary" onClick={this.toggle}>
              Cancel
            </Button>
          </ModalFooter>
        </Modal>
      </div>
    );
  }
}
