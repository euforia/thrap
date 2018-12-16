import React, { Component } from 'react';
import ReactJson from 'react-json-view';
import { Link } from 'react-router-dom';
import { Typography, Chip, Grid, Button, withStyles, Tooltip } from '@material-ui/core';

import KVList from '../common/KVList';

import {thrap} from '../api/thrap';

function getKVPairs(m) {
    if (m === null || m === undefined) return [];
    var pairs = Object.keys(m).map(k => {
        return {name: k, value: m[k], required: m[k] === '' ? true : false};
    })
    return pairs;
}

const styles = theme => ({
    heading: {
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

    render() {
        const { deploy } = this.state;
        const { classes } = this.props;
        const metaPairs = this.state.metas;
        const varsPairs = this.state.vars;
        const project = this.props.match.params.project;
       
        return (
            <div>
                <Grid container alignItems="center" justify="space-between">
                    <Grid item xs={2}>
                        <div style={{padding: '20px 0'}}>
                            <Typography variant="h5">{deploy.Name}</Typography>
                            <Typography>Profile: {deploy.Profile.ID}</Typography>
                        </div>
                    </Grid>
                    <Grid item xs={10} style={{textAlign:'right'}}>
                        <Tooltip title={deploy.StateMessage ? deploy.StateMessage : 'Status'}>
                            <Chip 
                                label={thrap.stateLabel(deploy.State,deploy.Status)} 
                                color={thrap.stateLabelColor(deploy.State, deploy.Status)}
                            />
                        </Tooltip>
                    </Grid>
                </Grid>
                {/* <Divider/> */}
                <div className={classes.btnPanel}>
                    <Button color="primary" variant="outlined"
                        component={Link} 
                        to={"/project/"+project+"/deploy/"+deploy.Profile.ID+"/"+deploy.Name+"/deploy"}
                    >
                        {thrap.stateLabel(deploy.State,deploy.Status)==='Deployed' ? 'Re-deploy' : 'Deploy'}
                    </Button>
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
            </div>
        );
    }
}

export default withStyles(styles)(Deployment);