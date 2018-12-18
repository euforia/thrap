import React, { Component } from 'react';
import ReactJson from 'react-json-view';
import { Link } from 'react-router-dom';
import { Typography, Chip, Grid, Button, withStyles, Tooltip } from '@material-ui/core';

import KVList from '../common/KVList';
import DeployStop from './DeployStop';

import {thrap} from '../api/thrap';

function getKVPairs(m) {
    if (m === null || m === undefined) return [];
    var pairs = Object.keys(m).map(k => {
        return {name: k, value: m[k], required: m[k] === '' ? true : false};
    })
    return pairs;
}

const styles = theme => ({
    header: {
        paddingBottom: theme.spacing.unit,
        paddingTop: theme.spacing.unit,
    },
    btnPanel: {
        textAlign: 'right',
        paddingBottom: theme.spacing.unit,
        paddingTop: theme.spacing.unit,
    },
    expPane: {
        border: 'none',
        boxShadow: 'none',
    }
})

class Deployment extends Component {

    constructor(props) {
        super(props);
        this.state = {
            deploy: {
                Profile: {
                    ID: '',
                },
            },
            profile: {},
            metas:[],
            vars:[],
            rawSpec: '',
            showStopModal: false,
        };

        this.fetchDeploy();
        this.fetchProfile();
    }
    
    fetchProfile() {
        const prof = this.props.match.params.profile;
        thrap.Profile(prof).then(resp => {
            this.setState({profile: resp.data});
        });
    }

    fetchDeploy() {
        const proj = this.props.match.params.project,
            prof = this.props.match.params.profile,
            inst = this.props.match.params.instance;

        thrap.Deployment(proj, prof, inst).then(resp => {
            var data = resp.data;
            var rawSpec = '';
            if (data.Spec !== undefined) {
                rawSpec = atob(data.Spec)
                data.Spec = JSON.parse(rawSpec);
            } else {
                data.Spec = {};
            }
            var metaPairs = getKVPairs(data.Profile.Meta !== undefined ? data.Profile.Meta : this.state.profile.Meta);
            var varsPairs = getKVPairs(data.Profile.Variables !== undefined ? data.Profile.Variables : this.state.profile.Variables);
            
            this.setState({
                rawSpec: rawSpec,
                vars: varsPairs,
                metas: metaPairs,
                deploy: data});
        });
    }

    // handleKVChange = pairs => {
    //     this.setState({vars: pairs});
    // }

    updateSpec = code => {
        console.log(code);
    }

    onClickStop = () => {
        const {project, profile, instance} = this.props.match.params;
        if (!thrap.isAuthd(profile)) {
            const path = `/login/${profile}#/project/${project}/deploy/${profile}/${instance}`;
            this.props.history.push(path);
            return;
        }
        this.setState({showStopModal:true});
    }

    hideDeployStopModal = () => {
        this.setState({showStopModal:false});
    }

    handleStop = (purge) => {
        const {project,profile,instance} = this.props.match.params;
        thrap.StopInstance(project, profile, instance, purge)
        .then(resp => {
            this.setState({showStopModal:false});
            this.fetchDeploy();
        })
        .catch(err => {
            console.log(err);
            this.setState({showStopModal:false});
            this.fetchDeploy();
        })    
    }

    render() {
        const { deploy } = this.state;
        const { classes } = this.props;
        const metaPairs = this.state.metas;
        const varsPairs = this.state.vars;
        const project = this.props.match.params.project;
        const status = thrap.stateLabel(deploy.State,deploy.Status);
       
        return (
            <div>
                <Grid container alignItems="center" className={classes.header}>
                    <Grid item xs={2}>
                        <Typography variant="h5">{deploy.Name}</Typography>
                        <Typography>Profile: {deploy.Profile.ID}</Typography>
                    </Grid>
                    <Grid item xs={10} style={{textAlign:'right'}}>
                        <Button color="secondary"
                            onClick={this.onClickStop}
                            disabled={status.includes('Deploy') ? false : true}
                        >
                            Stop
                        </Button>
                        <Button color="primary"
                            component={Link} 
                            to={"/project/"+project+"/deploy/"+deploy.Profile.ID+"/"+deploy.Name+"/deploy"}
                        >
                            {status==='Deployed' ? 'Re-deploy' : 'Deploy'}
                        </Button>
                    </Grid>
                </Grid>
                <div className={classes.btnPanel}>
                    <Tooltip title={deploy.StateMessage ? deploy.StateMessage : 'Status'} 
                        placement="left"
                    >
                        <Chip 
                            label={status} 
                            color={thrap.stateLabelColor(deploy.State, deploy.Status)}
                        />
                    </Tooltip>
                </div>
                <Grid container justify="space-between">
                    <Grid item xs={5}>
                        <Typography variant="h6" className={classes.heading}>Variables</Typography>
                        <KVList
                            readOnly={true}
                            pairs={varsPairs} 
                            // onKVChange={this.handleKVChange}
                        />
                        <Typography variant="h6" className={classes.heading}>Meta</Typography>
                        <KVList
                            readOnly={true}
                            pairs={metaPairs} 
                            // onKVChange={this.handleKVChange}
                        />
                    </Grid>
                    <Grid item xs={6}>
                        <Typography variant="h6" className={classes.heading}>Descriptor</Typography>
                        <ReactJson name="descriptor" 
                                src={this.state.deploy.Spec} 
                                style={{background: 'none', width:'100%', minHeight: '50px'}}
                                displayObjectSize={true} 
                                sortKeys={true} 
                                collapsed={3} 
                                // onEdit={this.onEdit}
                                displayDataTypes={false}
                                iconStyle="circle"
                                // theme="codeschool" 
                            />
                    </Grid>  
                </Grid>
                <DeployStop 
                    open={this.state.showStopModal}
                    onCancel={this.hideDeployStopModal}
                    onStop={this.handleStop}
                    name={deploy.Name}
                />
            </div>
        );
    }
}

export default withStyles(styles)(Deployment);