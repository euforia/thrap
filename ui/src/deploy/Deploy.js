import React, { Component } from 'react';
import { withStyles, Typography, Grid, Button } from '@material-ui/core';
import { Link } from 'react-router-dom';
import thrap from '../api/thrap';
import KVList from '../common/KVList';
import Descriptors from '../descriptor/Descriptors';
import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

const styles = theme => ({
    heading: {
        paddingTop: theme.spacing.unit,
        paddingBottom:theme.spacing.unit,
    },
    jsonViewer: {
        paddingTop: theme.spacing.unit*2,
        paddingBottom:theme.spacing.unit*2,
    },
    btnControls: {
        textAlign:'right',
    },
    modalCenter: {
        position: 'absolute',
        top: '50%',
        left: "50%",
        outline: "none",
        transform: 'translate(-50%, -50%)',
    }
})

function getKVPairs(m) {
    if (m === null || m === undefined) return [];
    var pairs = Object.keys(m).map(k => {
        return {name: k, value: m[k], required: m[k] === '' ? true : false};
    })
    return pairs;
}

class Deploy extends Component {
    constructor(props) {
        super(props);
        this.state = {
            vars: [],
            // metas: [],
            profile: {Variables:{}},
            errMsg: '',
            specName: '',
            specErr: false,
            disabled: false,
        };

        if (!thrap.isAuthd()) {
            const {project, profile, instance} = this.props.match.params;
            var path = `/login#/project/${project}/deploy/${profile}/${instance}/deploy`;
            this.props.history.push(path);
        } else {
            this.fetchProfile();
        }
    }

    fetchProfile() {
        const prof = this.props.match.params.profile;
        thrap.Profile(prof).then(resp => {
            var vs = getKVPairs(resp.data.Variables);
            // var ms = getKVPairs(resp.data.Meta);
            this.setState({
                profile: resp.data, 
                vars: vs
            });
        });
    }

    deploy = () => {
        // Check
        const {specName} = this.state;
        var specErr = false;
        if (specName === '') {
            specErr = true;
        }

        var vars = this.state.vars,
            payload = {},
            invalid = false;

        for (var i = 0; i < vars.length; i++) {
            var v = vars[i];
            if (v.value === '') {
                v.error = true;
                invalid = true;
                continue;
            }

            payload[v.name] = v.value;
        }

        this.setState({'vars': vars, 'specErr':specErr});
        if (invalid||specErr) {
            return;
        }

        // Start deploy call
        this.setState({disabled:true});
        const { project, profile, instance } = this.props.match.params;
        const req = {
            Vars: payload,
            Descriptor: specName,
        }
        thrap.DeployInstance(project, profile, instance, req)
            .then(data => {
                this.setState({disabled:false});
                // this.props.onDeploy();
                console.log(data);
            })
            .catch(error => {
                this.setState({disabled:false});
               
                var resp = error.response;
                this.setState({
                    errMsg: resp.data,
                });
                // this.props.onDeployError(resp.data);
            });
    }

    handleKVAdd = (key, value) => {
        var vars = this.state.vars;
        vars.push({name:key,value:value});
        this.setState({vars:vars});
    }

    handleKVChange = (i, event) => {
        var vars = this.state.vars;
        vars[i].value = event.currentTarget.value;
        this.setState({vars: vars});
    }

    handleKVRemove = (i) => {
        var vars = this.state.vars;
        vars.splice(i,1);
        this.setState({vars:vars});
    }

    onSpecSelected = (name) => {
        this.setState({specName: name});
    }
    
    render() {
        const { classes } = this.props;
        const { project, profile, instance } = this.props.match.params;
        const { disabled, specName, specErr } = this.state;
        // only pass in initial vars
        // const pairs = getKVPairs(this.state.profile.Variables);
        const pairs = this.state.vars;
        // const vars = this.state.profile.Variables;

        return (
            <div>
                <div style={{paddingTop: 20, paddingBottom:40}}>
                    <Grid container justify="space-between" alignItems="center">
                        <Grid item xs={8}>
                            <Typography variant="h6">
                                Deploy <b>{instance}</b> to <b>{profile}</b>
                            </Typography>
                        </Grid>
                        <Grid item xs={1} style={{textAlign:'right'}}>
                            <IconButton component={Link} 
                                to={'/project/'+project+'/deploy/'+profile+'/'+instance}
                            >
                                <CloseIcon/>
                            </IconButton>
                        </Grid>
                    </Grid>
                </div>
                <Typography color="secondary">{this.state.errMsg}</Typography>
                <Grid container justify="space-between">
                    <Grid item xs={5}>
                        <Typography variant="h6" className={classes.heading}>Variables</Typography>
                        <KVList
                            readOnly={disabled}
                            pairs={pairs}
                            onKVChange={this.handleKVChange} 
                            onKVAdd={this.handleKVAdd}
                            onKVRemove={this.handleKVRemove}
                        />
                        <div className={classes.btnControls}>
                            <Button color="secondary"
                                disabled={disabled}
                                component={Link} 
                                to={'/project/'+project+'/deploy/'+profile+'/'+instance}
                            >
                                Cancel
                            </Button>
                            <Button color="primary"
                                disabled={disabled}
                                onClick={this.deploy}
                            >
                                Deploy
                            </Button>
                        </div>
                    </Grid>
                    <Grid item xs={6}>
                        <Descriptors 
                            project={project} 
                            disabled={disabled}
                            spec={specName}
                            onChange={this.onSpecSelected}
                            error={specErr}/>
                    </Grid>
                </Grid>
            </div>
        );
    }
}

export default withStyles(styles)(Deploy);