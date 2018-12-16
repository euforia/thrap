import React, { Component } from 'react';
import { Grid, Button, FormControl, TextField, MenuItem, withStyles, Typography } from '@material-ui/core';
import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close'; 
import { Link } from 'react-router-dom';

import {thrap, validateName} from '../api/thrap';

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

        var profile = props.match.params.profile ? props.match.params.profile : '';
        this.state = {
            profile: profile,
            profileDisabled: profile !== '',
            instance: '',
            profErr: false,
            instErr: '',
            disabled: false,
            errMsg: '',
        };

        if (!thrap.isAuthd(profile)) {
            var path = profile === ''
                ? `/login#/project/${props.project}/deploys/new`
                : `/login/${profile}#/project/${props.project}/deploy/${profile}/new`;
            this.props.history.push(path);
        }
    }

    handleInstanceNameChange = (event) => {
        var val = event.target.value;
        
        this.setState({
            instance: val,
            instErr: validateName(val),
        })
    }

    handleProfileChange = (event) => {
        var val = event.target.value,
            props = this.props;

        // Make sure we are authd to the profile
        if (!thrap.isAuthd(val)) {
            var path = `/login/${val}#/project/${props.project}/deploy/${val}/new`;
            props.history.push(path);
            return;
        }

        this.setState({
            profile: val,
            profErr: false,
        })
    }

    handleCreateDeploy = event => {
        var s = {
            profErr: this.state.profile === '',
            instErr: validateName(this.state.instance),
        }

        if (s.profErr||s.instErr!=='') {
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
        const { profile, profErr, instance, instErr, disabled, profileDisabled } = this.state;
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
                        onChange={this.handleProfileChange}
                        margin="normal"
                        select
                        error={profErr}
                        fullWidth
                        required
                        disabled={disabled || profileDisabled}
                    >
                        {profiles.map(option => (
                            <MenuItem key={option.ID} value={option.ID}>
                                {option.Name}
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField 
                        label="Instance" 
                        value={instance}
                        onChange={this.handleInstanceNameChange}
                        required
                        margin="normal"
                        fullWidth
                        error={instErr!==''}
                        helperText={instErr!=='' ? instErr : 'Name of new instance'}
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