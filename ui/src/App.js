import React, { Component } from 'react';
// import { Route, Redirect } from 'react-router-dom';
// import { ToastContainer } from 'react-toastify';
// import 'react-toastify/dist/ReactToastify.css';
import axios from 'axios';
import constants from "./config"
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'

// import { Button, FormGroup, FormControl, ControlLabel } from "react-bootstrap";
import "./login.css";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      password: ""
    }
  }



  validateForm = () => {
    return this.state.name.length > 0 && this.state.password.length > 0;
  }

  componentDidMount() {
    if (localStorage.token) {
      if (localStorage.role !== "admin") {
        this.props.history.push('/details');
      } else {
        this.props.history.push('/home');
      }
    }
  }

  changeState = (event, param) => {
    this.setState({
      [param]: event
    })
  }

  handleSubmit = async (event) => {
    event.preventDefault();
    const response = await axios({
      method: "POST",
      url: constants.url + "login",
      data: {
        username: this.state.name,
        password: this.state.password
      },
      headers: {
        'content-type': 'application/json',
        Authorization: `Bearer `
      }
    });
    if (response.status === 200) {
      if (response.data) {
        localStorage.role = response.data.role
        localStorage.token = response.data.token
        if (response.data.role === "admin") {
          this.props.history.push('/home');
        } else {
          let [name, number] = localStorage.token.split(',?,')
          localStorage.uname = name
          localStorage.number = number
          this.props.history.push('/details');
        }
      }
    }
  }

  render() {
    return (
      <div className="Login">
        <Form onSubmit={this.handleSubmit}>
          <Form.Group controlId="formBasicEmail">
            <Form.Label>Name</Form.Label>
            <Form.Control value={this.state.name} onChange={e => this.changeState(e.target.value, 'name')} type="text" placeholder="Enter name" />
          </Form.Group>
          <Form.Group controlId="formBasicPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control type="text" value={this.state.password}
              onChange={e => this.changeState(e.target.value, 'password')} placeholder="Password" />
          </Form.Group>
          <Button variant="primary" disabled={!this.validateForm()} type="submit">
            Login
          </Button>
        </Form>
      </div>
    );
  }
}

export default App;
