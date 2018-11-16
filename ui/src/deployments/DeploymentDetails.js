import React, { Component } from 'react';

import { Grid } from '@material-ui/core';
import equal from 'fast-deep-equal';

import thrap from '../thrap.js';

const styles = ({
    headerDetails: {
        padding: '20px 0 40px 0',
        fontSize: '13px',
        borderBottom: '1px solid #668295',
    },
    section: {
        padding: '20px 0',
    },
    varkey: {
        opacity: 0.7,
    },
    varval: {
        minWidth: '160px',
        paddingLeft: '10px',
    },
    vars: {
        padding: '10px 0',
    },
    varsheader: {
        fontSize: '16px',
        padding: '10px 0',
        display: 'inline-block',
    },
    kvcontainer: {
        padding: '3px 0',
    },
});

class DeploymentDetails extends Component{
    constructor(props) {
        super(props);

        this.state = {
            environment: props.environment,
            deployment: props.deployment,
            messageClass: '',
        }

        var deploy = props.deployment;
        if (deploy.StateMessage !== undefined && deploy.StateMessage !== '') {
            // if (deploy.StateMessage.toLowerCase().includes('error')) {
            this.state.messageClass = 'error';
            // }
        }
    }

    getProfileMetaVars() {
        // Set default
        var env = this.state.environment;

        // Used the prepared version
        if (this.state.deployment.State > 1) {
            env = this.state.deployment.Profile;
        }
        console.log(env);

        return Object.keys(env.Meta).map(key => {
            var value = env.Meta[key];
            return (
                <div style={styles.kvcontainer} key={key}>
                    <Grid container spacing={0}>
                        <Grid item xs={6}><span style={styles.varkey}>{key}:</span></Grid>
                        <Grid item xs={6}>{value}</Grid>
                    </Grid>
                </div>
            );
        });
    }

    getProfileVars() {
        var env = this.state.environment;
        if (this.state.deployment.State > 1) {
            env = this.state.deployment.Profile;
        }

        return Object.keys(env.Variables).map(key => {
            var value = env.Variables[key];
            return (
                <div style={styles.kvcontainer} key={key}>
                    <Grid container spacing={0}>
                        <Grid item xs={6}><span style={styles.varkey}>{key}:</span></Grid>
                        <Grid item xs={6}>{value}</Grid>
                    </Grid>
                </div>
            );
        });
    }
    
    shouldComponentUpdate(nextProps) {
        return !equal(this.state.deployment, nextProps.deployment);
    }

    componentWillUpdate(nextProps) {
        this.setState({deployment:nextProps.deployment});
    }

    render() {
        var deploy = this.state.deployment,
            env = this.state.environment;

        var metaKV = this.getProfileMetaVars(),
            varsKV = this.getProfileVars();
        
        return (
            <Grid container spacing={0}>
                <Grid item xs={12}>
                    <div className="header-container">
                        <div className="header-title">{env.Name + " / " + deploy.Name}</div>
                        <div className="subscript">{"Previous: " +deploy.Previous}</div>
                    </div>
                </Grid>
                <Grid item xs={12}>
                    <table style={styles.headerDetails} className="kvtable">
                        <tbody>
                        <tr>
                            <td style={styles.varkey}>State:</td>
                            <td style={styles.varval}>{thrap.stateLabel(deploy.State, deploy.Status)}</td>
                            <td style={styles.varkey}>Message:</td>
                            <td style={styles.varval} className={this.state.messageClass}>{deploy.StateMessage}</td>
                        </tr>
                        <tr>
                            <td style={styles.varkey}>Version:</td>
                            <td style={styles.varval}>{deploy.Version}</td>
                            <td style={styles.varkey}>Nonce:</td>
                            <td style={styles.varval}>{deploy.Nonce}</td>
                        </tr>
                        <tr>
                            <td style={styles.varkey}>Created:</td>
                            <td style={styles.varval}>{(new Date(deploy.CreatedAt/1000000)).toLocaleString()}</td>
                            <td style={styles.varkey}>Modified:</td>
                            <td style={styles.varval}>{deploy.ModifiedAt === undefined ? '' : (new Date(deploy.ModifiedAt/1000000)).toLocaleString()}</td>
                        </tr>
                        </tbody>
                    </table>
                </Grid>
                <Grid item xs={6}>
                    <div style={styles.section}>
                        <div>
                            <div style={styles.varsheader} className="vars-header">Metadata</div>
                        </div>
                        <div style={styles.vars}>{metaKV}</div>
                    </div>
                </Grid>
                <Grid item xs={6}>
                    <div style={styles.section}>
                        <div>
                            <div style={styles.varsheader} className="vars-header">Variables</div>
                        </div>
                        <div style={styles.vars}>{varsKV}</div>
                    </div>
                </Grid>
            </Grid>
        );
    }
}

export default DeploymentDetails;