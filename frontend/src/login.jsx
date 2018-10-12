import React from "react";
import axios from "axios";
import { Button, Form, FormGroup, Label, Input } from "reactstrap";
import { Cookies } from "react-cookie";

export class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      username: "",
      password: "",
      invalid: this.props.invalid
    };
  }

  validateForm() {
    return this.state.username.length > 0 && this.state.password.length > 0;
  }

  handleChange = event => {
    this.setState({
      [event.target.id]: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();
    axios
      .post(process.env.REACT_APP_API_URL + "/login", null, {
        auth: { username: this.state.username, password: this.state.password }
      })
      .then(response => {
        var cookieValue = response.data;
        if (response.data === "") {
          cookieValue = "invalid";
        }
        new Cookies().set("X-Session-Token", cookieValue, { path: "/" });
        window.location.reload();
      });
  };

  render() {
    var loginErr;
    if (this.props.invalid) {
      loginErr = (
        <FormGroup>
          <Label id="error-text">Invalid username or password.</Label>
        </FormGroup>
      );
    }
    return (
      <div className="Login" id="login">
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
          {loginErr}
          <Button color="primary" disabled={!this.validateForm()}>
            Submit
          </Button>
        </Form>
      </div>
    );
  }
}
