import React from "react";
import axios from "axios";
import { Cookies } from "react-cookie";
import "./index.css";
import { Header } from "./header";
import { TaskList } from "./tasklist";
import { Login } from "./login";

export class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: null,
      user: { Username: "User", Email: "user@example.com", Tasks: null },
      invalid: false
    };
  }

  componentDidMount() {
    var token = new Cookies().get("X-Session-Token") || null;
    if (token != null) {
      axios
        .post("http://localhost:3000/token", null, {
          headers: { "X-Session-Token": token }
        })
        .then(response => {
          this.setState({ token: token, invalid: false });
          this.getUserFromToken(token);
        })
        .catch(error => {
          new Cookies().set("X-Session-Token", "invalid", { path: "/" });
          this.setState({ token: "invalid", invalid: true });
        });
    }
  }

  getUserFromToken(token) {
    axios
      .get("http://localhost:3000/users", {
        headers: { "X-Session-Token": token }
      })
      .then(response => {
        this.setState({ user: response.data });
      });
  }

  render() {
    if (this.state.token !== null && this.state.token !== "invalid") {
      return (
        <div>
          <Header token={this.state.token} user={this.state.user} />
          <div id="content">
            <TaskList token={this.state.token} user={this.state.user} />
          </div>
        </div>
      );
    } else {
      return (
        <div>
          <Header token={this.state.token} user={this.state.user} />
          <div id="content">
            <Login invalid={this.state.invalid} />
          </div>
        </div>
      );
    }
  }
}
