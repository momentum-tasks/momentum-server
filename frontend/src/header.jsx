import React from "react";
import axios from "axios";
import { Cookies } from "react-cookie";
import {
  Button,
  Collapse,
  DropdownToggle,
  DropdownMenu,
  DropdownItem,
  Form,
  FormGroup,
  Input,
  Label,
  Modal,
  ModalHeader,
  ModalBody,
  Navbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  UncontrolledDropdown,
  NavItem,
  NavLink
} from "reactstrap";

export class Header extends React.Component {
  constructor(props) {
    super(props);

    this.toggle = this.toggle.bind(this);
    this.state = {
      showCreateTaskModal: false,
      isOpen: false,

      name: "New Task",
      description: "",
      due: Date.now(),
      priority: 4294967295,
      completed: false
    };

    this.toggleShowCreateTaskModal = this.toggleShowCreateTaskModal.bind(this);
  }

  toggleShowCreateTaskModal() {
    this.setState({
      showCreateTaskModal: !this.state.showCreateTaskModal,

      name: "New Task",
      description: "",
      due: new Date().toISOString(),
      priority: 4294967295,
      completed: false
    });
  }

  validateCreateTaskModal() {
    return (
      this.state.name.length > 0 &&
      this.state.description.length > 0 &&
      this.state.due.toString().length > 0 &&
      this.state.priority > 0
    );
  }

  handleCreateTaskChange = event => {
    this.setState({
      [event.target.id]: event.target.value
    });
  };

  handleCreateTaskSubmit = event => {
    event.preventDefault();
    axios
      .post(
        process.env.REACT_APP_API_URL + "/tasks",
        {
          name: this.state.name,
          description: this.state.description,
          due: new Date(this.state.due).toISOString(),
          priority: this.state.priority,
          completed: this.state.completed
        },
        {
          headers: {
            "X-Session-Token": this.props.token,
            "Content-Type": "application/json"
          }
        }
      )
      .then(response => {
        window.location.reload();
      })
      .catch(error => {
        console.log(error.response);
      });
  };

  toggle() {
    this.setState({
      isOpen: !this.state.isOpen
    });
  }

  logout = event => {
    event.preventDefault();
    axios.post(process.env.REACT_APP_API_URL + "/logout", null, {
      headers: { "X-Session-Token": this.props.token }
    });
    new Cookies().remove("X-Session-Token");
    window.location.reload();
  };

  render() {
    var createTaskModal = (
      <Modal
        isOpen={this.state.showCreateTaskModal}
        toggle={this.toggleShowCreateTaskModal}
        className={this.props.className}>
        <ModalHeader toggle={this.toggleShowCreateTaskModal}>
          Create new task
        </ModalHeader>
        <ModalBody>
          <Form onSubmit={this.handleCreateTaskSubmit}>
            <FormGroup>
              <Label>Task Name</Label>
              <Input
                autoFocus={true}
                type="name"
                id="name"
                onChange={this.handleCreateTaskChange}
              />
            </FormGroup>
            <FormGroup>
              <Label>Description</Label>
              <Input
                type="textarea"
                id="description"
                onChange={this.handleCreateTaskChange}
              />
            </FormGroup>
            <FormGroup>
              <Label>Due Date</Label>
              <Input
                type="datetime-local"
                id="due"
                onChange={this.handleCreateTaskChange}
              />
            </FormGroup>
            <Button color="primary" disabled={!this.validateCreateTaskModal()}>
              Submit
            </Button>
          </Form>
        </ModalBody>
      </Modal>
    );

    return (
      <div>
        {createTaskModal}
        <Navbar color="dark" dark expand="md">
          <NavbarBrand href="/">Momentum Tasks</NavbarBrand>
          <NavbarToggler onClick={this.toggle} />
          <Collapse isOpen={this.state.isOpen} navbar>
            <Nav className="ml-auto" navbar>
              <NavItem onClick={this.toggleShowCreateTaskModal}>
                <NavLink>New Task</NavLink>
              </NavItem>
              <UncontrolledDropdown nav inNavbar>
                <DropdownToggle nav caret>
                  {this.props.user.Username}
                </DropdownToggle>
                <DropdownMenu right>
                  <DropdownItem>Settings</DropdownItem>
                  <DropdownItem divider />
                  <DropdownItem onClick={this.logout}>Logout</DropdownItem>
                </DropdownMenu>
              </UncontrolledDropdown>
            </Nav>
          </Collapse>
        </Navbar>
      </div>
    );
  }
}
