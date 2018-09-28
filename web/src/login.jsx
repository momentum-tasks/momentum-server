import React from "react";
import axios from "axios";
import { Button, Form, FormGroup, Label, Input } from "reactstrap";
import { Cookies } from "react-cookie";

export class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      username: "",
      password: ""
    };
  }

  validateForm() {
    return this.state.username.length > 0 && this.state.password.length > 0;
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
      .post("http://localhost:3000/login", null, {
        auth: { username: this.state.username, password: this.state.password }
      })
      .then(response => {
        new Cookies().set("X-Session-Token", response.data);
        window.location.reload();
      });
  };

  render() {
    return (
      <div className="Login">
        <Form onSubmit={this.handleSubmit}>
          <FormGroup>
            <Label>Username</Label>
            <Input
              autoFocus
              type="username"
              id="username"
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup>
            <Label>Password</Label>
            <Input type="password" id="password" onChange={this.handleChange} />
          </FormGroup>
          <Button color="primary" disabled={!this.validateForm()}>
            Submit
          </Button>
        </Form>
      </div>
    );
  }
}
