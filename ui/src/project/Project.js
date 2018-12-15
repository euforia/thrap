import React, { Component } from 'react';
import { Route } from 'react-router-dom';

import { Typography, Grid } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import { Link, Switch, Redirect } from 'react-router-dom';
import ConfirmDelete from '../common/ConfirmDelete';
import Deployments from '../deploy/Deployments';
import Deployment from '../deploy/Deployment';
import NewDeployment from '../deploy/NewDeployment';
import Deploy from '../deploy/Deploy';

import thrap from '../api/thrap';

const styles = theme => ({
    title: {
        paddingTop: theme.spacing.unit * 3,
        paddingBottom: theme.spacing.unit * 3,
        // paddingLeft: theme.spacing.unit * 4,
        // paddingRight: theme.spacing.unit * 4,
    },
    footer: {
        paddingTop: theme.spacing.unit * 2,
        paddingBottom: theme.spacing.unit * 2,
        textAlign: 'right',
    },
    sections: {
        // borderRight: '1px solid #ddd',
        borderLeft: '1px solid #ddd',
    }
});

class Project extends Component {
    constructor(props) {
        super(props);
        this.state = {
            project: {},
            showDelConfirm: false,
        };

        this.fetchProject();
    }

    fetchProject() {
        var id = this.props.match.params.project;
        thrap.Project(id).then(proj => {
            this.setState({project: proj.data});
        });
    }

    showDelConfirm = event => {
        this.setState({showDelConfirm: true});
    }

    hideDelConfirm = event => {
        this.setState({showDelConfirm: false});
    }
    
    handleDelete = event => {
        // API.DeleteProject(this.state.project.ID).then(resp => {
        //     this.props.history.push('/projects');
        // });
    }

    render() {
        const { classes } = this.props;
        const profiles = this.props.profiles;
        const project = this.state.project;

        return (
            <div>
                <Typography variant="h4" className={classes.title}>
                    <Link to={'/project/'+project.ID}>{project.Name}</Link>
                </Typography>
                <Grid container>
                    <Grid item xs={12}>
                        <Switch>
                            <Route exact 
                                path="/project/:project/deploys" 
                                render={(props) => <Deployments {...props} project={project.ID} />} 
                            />
                            <Route exact
                                path="/project/:project/deploys/new" 
                                render={(props) => <NewDeployment {...props} project={project.ID} profiles={profiles} /> } 
                            />
                            <Route exact 
                                path="/project/:project/deploy/:profile/:instance" 
                                render={(props) => <Deployment {...props} />} 
                            />
                            <Route exact 
                                path="/project/:project/deploy/:profile/:instance/deploy" 
                                render={(props) => <Deploy {...props} />} 
                            />
                            <Redirect to="/project/:project/deploys" />
                        </Switch>
                    </Grid>
                </Grid>
                <div className={classes.footer}>
                    {/* <Button color="secondary" onClick={this.showDelConfirm}>Delete</Button> */}
                </div>
                <ConfirmDelete 
                    entity={project}
                    entityType='Project'
                    open={this.state.showDelConfirm} 
                    onCancel={this.hideDelConfirm} 
                    onDelete={this.handleDelete}
                />

            </div>
        );
    }
}

export default withStyles(styles)(Project);