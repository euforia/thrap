import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import Button from '@material-ui/core/Button';
import { Paper, FormControl, TextField, MenuItem } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import { Grid } from '@material-ui/core';
// import API from '../api/api';
import thrap from '../api/thrap';

const styles = theme => ({
  paper: {
    padding: theme.spacing.unit * 2,
    width: 400,
    margin: "auto",
    marginTop: theme.spacing.unit*8,
  },
  footer: {
    textAlign: 'right',
  },
  errorMessage: {
    textAlign: 'center',
    fontWeight: 'bold'
  }
});


class Login extends Component {
    constructor(props) {
        super(props);

        this.state = {
            profile: 'dev',
            providers: {},
            provider: 'vault',
            authType: 'token',
            username: '',
            password: '',
            token: '',
            userErr: false,
            passErr: false,
            tokenErr: false,
            errorMessage: '',
        }

        this.fetchAuthProviders();
        if (thrap.isAuthd()) {
            props.onLogin();
        }
    }

    fetchAuthProviders() {
        thrap.AuthMethods().then(methods => {
            var auths = {};
            for(var i = 0; i< methods.length; i++) {
                auths[methods[i].id] = methods[i];
            }
            this.setState({providers: auths});
        });
    }

    handleChange = name => event => {
        this.setState({
          [name]: event.target.value,
        });
    };

    handleProviderChange = event => {
        var val = event.target.value;
        var method = this.state.providers[val];
        this.setState({
            provider: val,
            authType: method.type,
        });
    }

    onLogin = (event) => {
        var s, err,
            payload = {
                provider: this.state.provider,
                type: this.state.authType,
            };
        
        if (this.state.authType === 'userpass') {
            s = {
                userErr: this.state.username === '',
                passErr: this.state.password === '',
            };
            if (s.userErr || s.passErr) err = true;
            payload['username'] = this.state.username;
            payload['password'] = this.state.password;    
        } else {
            s = {
                tokenErr: this.state.token === '',
            };
            if (s.tokenErr) err = true;
            payload['token'] = this.state.token;
        }

        this.setState(s);
        if (err) return; 

        thrap.Authenticate(this.state.profile, this.state.token)
        .then(resp => {
            this.props.onLogin({}, event);
        })
        .catch(error => {
            this.setState({
                errorMessage: 'Invalid password.  Please try again.',
                passErr: true,
            });
        });    
    }

    render() {
        const { classes } = this.props;
        const profiles = this.props.profiles;
        const providers = this.state.providers;

        return (
            <Grid container>
                <Grid item xs={12}>
                    <Paper className={classes.paper}>
                        <Typography margin="normal" color="secondary" className={classes.errorMessage}>
                            {this.state.errorMessage}
                        </Typography>
                        <FormControl fullWidth margin="normal">
                            <TextField label="Profile"
                                value={this.state.profile}
                                onChange={this.handleChange('profile')}
                                variant="outlined"
                                margin="normal"
                                select
                                fullWidth
                                required
                            >
                                {profiles.map(option => (
                                    <MenuItem key={option.ID} value={option.ID}>
                                        {option.Name}
                                    </MenuItem>
                                ))}
                            </TextField>
                            <TextField label="Method"
                                value={this.state.provider}
                                onChange={this.handleProviderChange}
                                variant="outlined"
                                margin="normal"
                                select
                                fullWidth
                                required
                            >
                                {Object.keys(providers).map(key => (
                                    <MenuItem key={providers[key].id} value={providers[key].id}>
                                        {providers[key].name}
                                    </MenuItem>
                                ))}
                            </TextField>
                            {this.state.authType === 'userpass' ?
                            <div>
                                <TextField label="Username"
                                    value={this.state.username}
                                    onChange={this.handleChange('username')}
                                    variant="outlined"
                                    margin="normal"
                                    error={this.state.userErr}
                                    fullWidth
                                    required
                                />
                                <TextField label="Password"
                                    type="password"
                                    value={this.state.password}
                                    onChange={this.handleChange('password')}
                                    variant="outlined"
                                    margin="normal"
                                    error={this.state.passErr}
                                    fullWidth
                                    required
                                /> 
                            </div> :
                            <TextField label="Token"
                                type="password"
                                value={this.state.token}
                                onChange={this.handleChange('token')}
                                variant="outlined"
                                margin="normal"
                                error={this.state.tokenErr}
                                fullWidth
                                required
                            />}
                        </FormControl>
                        <div className={classes.footer}>
                            <Button color="primary" onClick={this.onLogin}>
                                Sign In
                            </Button>
                        </div>
                    </Paper>
                </Grid>
            </Grid>
        );
    }
}

Login.propTypes = {
    classes: PropTypes.object,
};
  
export default withStyles(styles)(Login);