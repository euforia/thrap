import React, { Component } from 'react';

import { withStyles } from '@material-ui/core/styles';
import { Button } from '@material-ui/core';
import FormControl from '@material-ui/core/FormControl';
import TextField from '@material-ui/core/TextField';

import thrap from './thrap.js';

import './Login.css';

const styles = theme => ({
    loginButton: {
        margin: '10px 20px',
    },
    formControl: {
        width: 300,
    },
    helperText: {
        color: '#abbbc6'
    },
    selectInput: {
        textAlign: 'left'
    },
});

class Login extends Component {
    constructor(props) {
        super(props);

        this.handleSelectChange = this.handleSelectChange.bind(this);
        this.onPasswordChange = this.onPasswordChange.bind(this);
        
        this.onLogin = this.onLogin.bind(this);

        this.state = {
            password: '',
            invalidPassword: false,
            profile: '',
            profiles: [],
            errorMessage: '',
        };

        if (thrap.isAuthd()) {
            props.onLoginSucceeded();
        } else {
            // Fetch auth profiles for user to select
            this.fetchProfiles();
        }
    }

    fetchProfiles() {
        thrap.environments()
            .then(({data}) => {
                // console.log(data);
                this.setState({
                    profiles: data,
                    profile: data[0].ID,
                });
            })
            .catch(error => {
                console.log(error);
                console.log(error.repsonse);
            });
    }

    handleSelectChange = name => event => {
        this.setState({ 
            [name]: event.target.value
        });
    };

    onPasswordChange(event) {
        var s = {
            password: event.target.value,
        };

        if (s.password.length > 0) {
            s.invalidPassword = false;
        }

        this.setState(s);
    }

    onLogin() {
        console.log("TODO: Try login");

        if (this.state.password === "") {
            this.setState({invalidPassword: true});    
            return;
        }

        // Reset invalid password state
        this.setState({invalidPassword: false});

        thrap.authenticate(this.state.profile, this.state.password)
            .then(({data}) => {
                // console.log(data);
                // Call parent login callback
                this.props.onLoginSucceeded();
            })
            .catch(error => {
                this.setState({
                    errorMessage: 'Invalid password.  Please try again.',
                    invalidPassword: true,
                });
            });
    }

    
    render() {
        const { classes } = this.props;
        const menuItems = (
            this.state.profiles.map(p => (
                <option key={p.ID} value={p.ID}>{p.Name}</option>
            ))
        );

        return (
            <div style={{textAlign: 'center', padding: '100px'}} className={classes.helperText}>
                <div id="login-logo"></div>
                <div style={{padding: '20px 10px'}} className="error-container">{this.state.errorMessage}</div>
                <FormControl className={classes.formControl}>
                    <TextField
                        id="profile-input"
                        label="Profile"
                        // InputLabelProps={{className:classes.inputLabel}}
                        FormHelperTextProps={{className:classes.helperText}}
                        InputProps={{className:classes.selectInput}}
                        select
                        SelectProps={{
                            native: true,
                        }}
                        value={this.state.profile}
                        onChange={this.handleSelectChange('profile')}
                        helperText="profile to authenticate against"
                        margin="normal"
                        required
                    >
                        {menuItems}
                    </TextField>
                    <TextField
                        id="password-input"
                        label="Token"
                        // InputLabelProps={{className:classes.inputLabel}}
                        FormHelperTextProps={{className:classes.helperText}}
                        // InputProps={{className:classes.input}}
                        type="password"
                        autoComplete="current-password"
                        helperText="Vault token associated to the profile"
                        margin="normal"
                        value={this.state.password}
                        onChange={this.onPasswordChange}
                        error={this.state.invalidPassword}
                    />
                </FormControl>
                <div style={{marginTop:'20px'}}>
                    <Button variant="contained" color="primary" className={classes.loginButton} onClick={this.onLogin}>Login</Button>
                </div>

            </div>
        );
    }
}

// export default Login;
export default withStyles(styles)(Login);
