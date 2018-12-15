import React, { Component } from 'react';
import { Grid, Button, FormControl, TextField, MenuItem, withStyles, Typography } from '@material-ui/core';
import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close'; 
import { Link } from 'react-router-dom';

import thrap from '../api/thrap';

const styles = theme => ({
    footer: {
        textAlign: 'right',
        paddingTop: theme.spacing.unit * 2,
        paddingBottom: theme.spacing.unit * 2,
    },
    title: {
        paddingTop: theme.spacing.unit*2,
        paddingBottom: theme.spacing.unit*2,
    }
});

class NewDeployment extends Component {
    constructor(props) {
        super(props);
        this.state = {
            profile: '',
            instance: '',
            profErr: false,
            instErr: false,
            disabled: false,
            errMsg: '',
        };

        if (!thrap.isAuthd()) {
            this.props.history.push(`/login#/project/${props.project}/deploys/new`);
        }
    }

    handleChange = name => event => {
        var e = name.substr(0,4)+'Err',
            val = event.target.value;
        this.setState({
            [name]: val,
            [e]: false,
        });
    }

    handleCreateDeploy = event => {
        var s = {
            profErr: this.state.profile === '',
            instErr: this.state.instance === '',
        }

        if (s.profErr||s.instErr) {
            this.setState(s);
            return;
        }

        s.disabled = true;
        this.setState(s);
        
        thrap.CreateDeployment(this.props.project, this.state.profile, this.state.instance)
        .then(resp => {
            this.setState({disabled:false});
            var path = '/project/'+this.props.project+'/deploy/'+this.state.profile+'/'+this.state.instance;
            this.props.history.push(path);
        })
        .catch(err => {
            if (err.response) {
                this.setState({
                    errMsg: err.response.data,
                    disabled:false,
                });
                return;
            }
            console.log(err);
        });
    }

    render() {
        const { project, profiles } = this.props;
        const { profile, disabled } = this.state;
        const { classes } = this.props;
        
        return (
            <div>
                <Grid container alignItems="center" justify="space-between">
                    <Grid item xs={9}>
                        <Typography variant="h5" className={classes.title}>
                            New Deployment
                        </Typography>
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'right'}}>
                        <IconButton component={Link} to={'/project/'+project+'/deploys'}>
                            <CloseIcon/>
                        </IconButton>
                    </Grid>
                </Grid>
                <Typography color="secondary" className={classes.errMsg}>
                    {this.state.errMsg}
                </Typography>
                <FormControl fullWidth>
                    <TextField label="Profile"
                        value={profile}
                        onChange={this.handleChange('profile')}
                        margin="normal"
                        select
                        error={this.state.profErr}
                        fullWidth
                        required
                        disabled={disabled}
                    >
                        {profiles.map(option => (
                            <MenuItem key={option.ID} value={option.ID}>
                                {option.Name}
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField 
                        label="Instance" 
                        value={this.state.instance}
                        onChange={this.handleChange('instance')}
                        helperText=''
                        required
                        margin="normal"
                        fullWidth
                        error={this.state.instErr}
                        disabled={disabled}
                    />
                </FormControl>
                <div className={classes.footer}>
                    <Button component={Link} color="secondary"
                        to={'/project/'+project+'/deploys'}
                        disabled={disabled}
                    >
                        Cancel
                    </Button>
                    <Button color="primary" onClick={this.handleCreateDeploy}
                        disabled={disabled}
                    >
                        Create
                    </Button>
                </div>
            </div>
        );
    }
}

export default withStyles(styles)(NewDeployment);