import React, { Component } from 'react';
import axios from 'axios';
import constants from ".././config"
import Card from "react-bootstrap/Card"
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Container from 'react-bootstrap/Container'
import CardDeck from 'react-bootstrap/CardDeck'
import Modal from 'react-bootstrap/Modal'
import Form from 'react-bootstrap/Form'

class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            users: localStorage.users ? JSON.parse(localStorage.users) : [],
            show: false,
            user: {},
            add: false,
            edit:false
        }
    }

    closePopup = () => {
        this.setState({ show: false,edit:false })
    }

    closeAddPopup = () => {
        this.setState({ add: false })
    }

    getUsers = async () => {
        const response = await axios({
            method: "GET",
            url: constants.url + "users",
            headers: {
                'content-type': 'application/json',
                Authorization: `Bearer ${localStorage.token}`
            }
        });
        if (response.status === 200) {
            localStorage.users = JSON.stringify(response.data.users)
            if (response.data) {
                this.setState({
                    users: response.data.users
                })
            }
        }
    }


    componentDidMount() {
        localStorage.user = JSON.stringify({})
        if (localStorage.token) {
            if (localStorage.role !== "admin") {
                this.props.history.push('/details');
            } else {
                this.getUsers()
            }
        } else {
            this.props.history.push('/')
        }
    }

    changeState(event, param) {
        this.setState({
            [param]: event
        })
    }

    changeD = (b) => {
        localStorage.user = JSON.stringify(b)
        localStorage.uname = b.name
        localStorage.number = b.number
        this.props.history.push('/details')
    }

    logOut = () => {
        localStorage.clear()
        this.props.history.push('/')
    }

    Add = () => {
        this.setState({ show: true, user: { name: "", number: "", amount: 0, transactions: [] } })
    }

    Edit = (b) => {
        this.setState({ show: true, edit: true, user: { id: b.id, name: b.name, number: b.number, amount: b.amount, newval: 0, desc: "", transactions: b.transactions } })
    }


    Update = (b) => {
        this.setState({ add: true, user: { id: b.id, name: b.name, number: b.number, amount: b.amount, newval: 0, desc: "", transactions: b.transactions } })
    }

    save = async (event) => {
        event.preventDefault();
        let user = JSON.parse(JSON.stringify(this.state.user))
        user.newval = Number(user.newval)
        if (user.newval !== 0) {
            user.amount = user.amount + user.newval
            user.transactions = []
            user.transactions.push({ time: new Date().getTime(), am: user.newval, desc: user.desc })
        }
        const response = await axios({
            method: "PUT",
            url: constants.url + "users/" + user.id,
            data: user,
            headers: {
                'content-type': 'application/json',
                Authorization: `Bearer ${localStorage.token}`
            }
        });
        if (response.status === 200) {
            this.getUsers()
            this.setState({ add: false,edit: false })
        }
    }


    changeUserState = (event, param) => {
        let user = JSON.parse(JSON.stringify(this.state.user))
        if (param === "amount") event = Number(event)
        user[param] = event
        this.setState({
            user: user
        })
    }

    validateForm = () => {
        if (!this.state.user.name) return false
        return this.state.user.name.length > 0 && this.state.user.number.length > 0;
    }


    handleSubmit = async (event) => {
        event.preventDefault();
        let user = JSON.parse(JSON.stringify(this.state.user))
        user.amount = Number(user.amount)
        if (user.amount !== 0) {
            if (!user.transactions) user.transactions = []
            user.transactions.push({ time: new Date().getTime(), am: user.amount, desc: "" })
        }
        const response = await axios({
            method: "POST",
            url: constants.url + "users",
            data: user,
            headers: {
                'content-type': 'application/json',
                Authorization: `Bearer ${localStorage.token}`
            }
        });
        if (response.status === 201) {
            if (response.data) {
                let users = JSON.parse(JSON.stringify(this.state.users))
                users.push(response.data.user)
                localStorage.users = JSON.stringify(users)
                this.setState({ users: users, show: false })
            }
        }
    }

    render() {
        return (
            <div className="">

                <div style={{ marginTop: "10px" }}>
                    <Container>
                        <div>
                            <Button variant="danger" onClick={this.logOut}>Log Out</Button>
                            <Button variant="primary" onClick={this.Add}>Add User</Button>
                        </div>
                        <Row style={{ paddingTop: "10px" }}>
                            <CardDeck>
                                {this.state.users.map((board, i) =>
                                    <Col xs={12} md={4} lg={4} key={i} style={{ padding: "5px 0px" }}>
                                        <Card style={{ width: '18rem', 'backgroundColor': 'darkcyan', color: 'white' }}>
                                            <Card.Body>
                                                <Card.Title>{board.name}, {board.number}</Card.Title>
                                                <Card.Text>
                                                    Amount: {board.amount}
                                                </Card.Text>
                                                <Button onClick={() => this.Edit(board)} variant="secondary">Edit</Button>
                                                <Button onClick={() => this.Update(board)} variant="primary">Update</Button>
                                                <Button onClick={() => this.changeD(board)} variant="info">Transactions</Button>
                                            </Card.Body>
                                        </Card>
                                    </Col>
                                )}
                            </CardDeck>
                        </Row>
                    </Container>
                </div>



                <Modal
                    show={this.state.show}
                    onHide={() => this.closePopup(false)}
                    size="lg"
                    aria-labelledby="contained-modal-title-vcenter"
                    centered
                >
                    <Modal.Header closeButton>
                        <Modal.Title id="contained-modal-title-vcenter">
                            {this.state.edit ? 'Edit':'Add'} User
                        </Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <Form onSubmit={() => {
                            return this.state.edit ? this.save() : this.handleSubmit()
                            }}>
                            <Form.Group controlId="formBasicEmail">
                                <Form.Label>Name</Form.Label>
                                <Form.Control value={this.state.user.name} onChange={e => this.changeUserState(e.target.value, 'name')} type="text" placeholder="Enter name" />
                            </Form.Group>
                            <Form.Group controlId="formBasicPassword">
                                <Form.Label>Number</Form.Label>
                                <Form.Control type="text" value={this.state.user.number}
                                    onChange={e => this.changeUserState(e.target.value, 'number')} placeholder="Number" />
                            </Form.Group>
                            <Form.Group controlId="formBasicAmount">
                                <Form.Label>Amount</Form.Label>
                                <Form.Control type="text" value={this.state.user.amount}
                                    onChange={e => this.changeUserState(e.target.value, 'amount')} placeholder="amount" />
                            </Form.Group>
                            <Button variant="primary" disabled={!this.validateForm()} type="submit">
                                Save
                             </Button>
                        </Form>
                    </Modal.Body>
                    <Modal.Footer>
                        <Button onClick={this.closePopup}>Close</Button>
                    </Modal.Footer>
                </Modal>


                <Modal
                    show={this.state.add}
                    onHide={this.closeAddPopup}
                    size="sm"
                    aria-labelledby="contained-modal-title-vcenter"
                    centered
                >
                    <Modal.Header closeButton>
                        <Modal.Title id="contained-modal-title-vcenter">

                        </Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <Form onSubmit={this.save}>

                            <Form.Group controlId="formBasicPassword">
                                <Form.Label>Description</Form.Label>
                                <Form.Control type="text" value={this.state.user.desc}
                                    onChange={e => this.changeUserState(e.target.value, 'desc')} placeholder="Number" />
                            </Form.Group>
                            <Form.Group controlId="formBasicAmount">
                                <Form.Label>Amount</Form.Label>
                                <Form.Control type="text" value={this.state.user.newval}
                                    onChange={e => this.changeUserState(e.target.value, 'newval')} placeholder="amount" />
                            </Form.Group>
                            <Button variant="primary" type="submit">
                                Save
                             </Button>
                        </Form>
                    </Modal.Body>
                    <Modal.Footer>
                        <Button onClick={this.closeAddPopup}>Close</Button>
                    </Modal.Footer>
                </Modal>



            </div>




        );
    }
}

export default Home;
