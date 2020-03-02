import React, { Component } from 'react';
import axios from 'axios';
import constants from ".././config"
import Card from "react-bootstrap/Card"
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
// import Col from 'react-bootstrap/Col'
import Container from 'react-bootstrap/Container'
import ListGroup from 'react-bootstrap/ListGroup'
import Accordion from 'react-bootstrap/Accordion'

class Details extends Component {
    constructor(props) {
        super(props);
        this.state = {
            user: localStorage.user ? JSON.parse(localStorage.user) : {}
        }
    }

    getUsers = async () => {
        let [name, number] = [localStorage.uname, localStorage.number]
        const response = await axios({
            method: "GET",
            url: constants.url + "users/" + name + "/" + number,
            headers: {
                'content-type': 'application/json',
                Authorization: `Bearer ${localStorage.token}`
            }
        });
        if (response.status === 200) {
            if (response.data) {
                this.setState({
                    user: response.data.user
                })
                localStorage.user = JSON.stringify(response.data.user)
            }
        }
    }


    componentDidMount() {
        if (localStorage.token) {
            this.getUsers()
        }
    }

    changeState = (event, param) => {
        this.setState({
            [param]: event
        })
    }

    logOut = () => {
        localStorage.clear()
        this.props.history.push('/')
    }

    goBack = () => {
        this.props.history.push('/home')

    }


    render() {
        return (
            <div className="">

                <Container>
                    <Row>
                        <div>
                            {localStorage.role === "admin" ? (<Button variant="info" onClick={this.goBack}>Back</Button>) : ""}
                            {localStorage.role !== "admin" ? (<Button variant="danger" onClick={this.logOut}>Log Out</Button>) : ""}
                        </div>
                    </Row>
                    <Row>
                        {this.state.user.name ? (<Card style={{ width: '18rem' }}>
                            <Card.Body>
                                <Card.Title>{this.state.user.name}</Card.Title>
                                <Card.Text>
                                    Amount: {this.state.user.amount}
                                </Card.Text>
                            </Card.Body>
                            <Accordion defaultActiveKey="0">
                                <Card.Header>
                                    <Accordion.Toggle as={Button} variant="flush" eventKey="0">
                                        Transactions
                              </Accordion.Toggle>
                                </Card.Header>
                                <Accordion.Collapse eventKey="0">
                                    <ListGroup variant="flush">
                                        {this.state.user.transactions.map((board, i) =>
                                            <ListGroup.Item key={i}>
                                                Amount: {board.am} , Desc : {board.desc}, date : {new Date(board.time).toLocaleString()}
                                            </ListGroup.Item>
                                        )}
                                    </ListGroup>
                                </Accordion.Collapse>
                            </Accordion>
                            {/* <Card.Header>Transactions</Card.Header>
                            <ListGroup variant="flush">
                            {this.state.user.transactions.map((board, i) =>
                                <ListGroup.Item key={i}></ListGroup.Item>
                            )}
                            </ListGroup> */}
                        </Card>) : "No Data"}
                    </Row>
                </Container>

            </div>
        );
    }
}

export default Details;
