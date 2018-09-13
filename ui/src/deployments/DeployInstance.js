import React, { Component } from 'react';

import { Button } from '@material-ui/core';

import KeyValuePairs from '../common/KeyValuePairs.js';
import thrap from '../thrap.js';

const styles = ({
    container: {
        padding: '20px 40px',
        borderRadius: '5px',
        maxWidth: '500px',
    }
});

// Actual deploy component
class DeployInstance extends Component{
    constructor(props) {
        super(props);

        this.deployInstance = this.deployInstance.bind(this);
        this.onKVChange = this.onKVChange.bind(this);

        this.state = {
            project: props.project,
            environment: props.environment,
            deployment: props.deployment,
            variables: [],
            status: '',
        }

        this.state.variables = Object.keys(props.environment.Variables).map(key => {
            var val = props.environment.Variables[key];
            return {
                name: key,
                value: val,
                // for KV component remove button
                required: (val === ''),
                valueClass: '',
                keyClass: '',
            }
        });
    }

    onKVChange(pairs) {
        this.setState({
            variables: pairs,
        })
    }

    deployInstance() {
        var project = this.state.project.ID,
            env = this.state.environment.ID,
            deploy = this.state.deployment.Name,
            vars = this.state.variables;

        var payload = {},
            invalid = false;

        for (var i = 0; i < vars.length; i++) {
            var v = vars[i];
            if (v.value === '') {
                v.valueClass = 'invalid-input';
                invalid = true;
            }
            if (v.name === '') {
                v.keyClass = 'invalid-input';
                invalid = true;
            }

            payload[v.name] = v.value;
        }

        if (invalid) {
            this.setState({variables: vars});
            return;
        }

        thrap.deployInstance(project, env, deploy, payload)
            .then(data => {
                this.props.onDeploy();
            })
            .catch(error => {
                var resp = error.response;
                this.setState({
                    status: resp.data,
                });
                this.props.onDeployError(resp.data);
            });
    }

    render() {
        var deploy = this.state.deployment,
            env = this.state.environment;

        return (
            <div style={styles.container} className="theme-bg">
                <div className="header-container theme-color" style={{textAlign:'center'}}>
                    <div className="header-title">Deploy: {env.ID} / {deploy.Name}</div>
                </div>
                <div className="error-container">{this.state.status}</div>
                <KeyValuePairs title="Variables" pairs={this.state.variables} onKVChange={this.onKVChange}/>      
                <div className="create-btn-container" style={{textAlign:'center'}}>
                    <Button variant="contained" color="primary" onClick={this.deployInstance}>Deploy</Button>
                </div>
            </div>
        );
    }
}

export default DeployInstance;