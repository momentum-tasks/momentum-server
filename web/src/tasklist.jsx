import React from "react";
import axios from "axios";
import { Task } from "./task";

export class TaskList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      restData: []
    };
  }

  componentDidMount() {
    axios
      .get("http://localhost:3000/tasks", {
        headers: { "X-Session-Token": this.props.token }
      })
      .then(response => {
        this.setState({ restData: response.data });
      });
  }

  render() {
    var tasks = [];
    this.state.restData.forEach(function(item) {
      tasks.push(<Task key={item.ID} task={item} />);
    });
    return <div className="task-list">{tasks}</div>;
  }
}
