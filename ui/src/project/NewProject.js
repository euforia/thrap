import React, { Component } from 'react';
import { Grid, FormControl, TextField, Button, withStyles, Typography, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

import { Link } from 'react-router-dom';

import thrap from '../api/thrap';

const styles = theme => ({
  footer: {
    padding: theme.spacing.unit * 2,
    textAlign: 'right',
  },
  title: {
    paddingTop: theme.spacing.unit * 3,
    paddingBottom: theme.spacing.unit * 3,
  },
  errMsg: {
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit,
  }
});

class NewProject extends Component {
    constructor(props) {
        super(props);
        this.state = {
            id: '',
            name: '',
            nameErr: false,
            errMsg: '',
            disabled: false,
        }
        if (!thrap.isAuthd()) {
            this.props.history.push(`/login#/projects/new`);
        }
    }

    handleChange = name => event => {
        var val = event.target.value;
        this.setState({
            [name]: val,
            nameErr: val.includes(" "),
            id: val.toLowerCase(),
        })
    }

    handleCreateProject = event => {
        var {name}  = this.state;
        var s = {
            nameErr: name === ''
        }
        this.setState(s);
        if (s.nameErr) return;


        this.setState({disabled:true});
        // var proj = this.state;
        var obj = {
            Project: {
                ID: name,
            }
        };
        thrap.CreateProject(obj).then(resp => {
            this.setState({disabled:false});

            this.props.history.push("/project/"+name);
        })
        .catch(error => {
            var resp = error.response;
            this.setState({
                disabled: false,
                errMsg: resp.data,
            });
        });
    }

    render() {
        const { classes } = this.props;
        const { disabled } = this.state;
        return (
            <div>
                <Grid container alignItems="center" justify="space-between">
                    <Grid item xs={9}>
                        <Typography variant="h5" className={classes.title}>
                            New Project
                        </Typography>
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'right'}}>
                        <IconButton component={Link} to="/projects">
                            <CloseIcon/>
                        </IconButton>
                    </Grid>
                </Grid>
                <Typography color="secondary" className={classes.errMsg}>
                    {this.state.errMsg}
                </Typography>
                <FormControl fullWidth>
                    <TextField 
                        label="ID" 
                        value={this.state.id}
                        margin="normal"
                        fullWidth
                        disabled
                    />
                    <TextField 
                        label="Name" 
                        value={this.state.name}
                        onChange={this.handleChange('name')}
                        required
                        margin="normal"
                        fullWidth
                        error={this.state.nameErr}
                        disabled={disabled}
                    />
                </FormControl>
                <div className={classes.footer}>
                    <Button component={Link} to="/projects"
                        disabled={disabled}
                        color="secondary"
                    >
                        Cancel
                    </Button>
                    <Button color="primary" disabled={disabled}
                        onClick={this.handleCreateProject}>Create</Button>
                </div>
            </div>
        );
    }
}

export default withStyles(styles)(NewProject);