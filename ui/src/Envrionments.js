import React, { Component } from 'react';
import { Add } from '@material-ui/icons';
// import CreateEnvironment from './CreateEnvironment.js';

class Environments extends Component {
    constructor(props) {
        super(props);
    
        this.onEnvCreate = this.onEnvCreate.bind(this);
        this.onEnvCreated = this.onEnvCreated.bind(this);
        this.onCloseDialogue = this.onCloseDialogue.bind(this);

        this.state = {
            filter: '',
            createEnv: false,
            environments: this.getEnvs(),
        }
    }

    getEnvs() {
        return [{
            ID: "dev",
            Name: "Development",
            Configs: [{
                Name: "Orchestrator",
                Value: "http://nomad.service.consul:4646",
            },{
                Name: "Logs",
                Value: "http://logs.service.consul",
            },{
                Name: "Metrics",
                Value: "http://metrics.service.consul",
            },{
                Name: "Alerts",
                Value: "http://alerts.service.consul",
            },{
                Name: "Secrets",
                Value: "http://vault.service.consul:8200",
            }]
        },{
            ID: "int",
            Name: "Integration",
            Configs: [],
        }, {
            ID: "prod",
            Name: "Production",
            Configs: [],
        }];
    }

    onEnvCreate() {
        this.setState({
            createEnv: true,
        });
    }

    onEnvCreated(name) {
        console.log(name)
        this.setState({
            createEnv: false,
        });
    }
    
    onCloseDialogue() {
        this.setState({
            createEnv: false,
        });
    }

    render() {

        // if (this.state.createEnv) {
        //     return (
        //         <CreateEnvironment onEnvCreated={this.onEnvCreated} onCloseDialogue={this.onCloseDialogue}/>
        //     );
        // }

        var items = [];
        var envs = this.state.environments
        for (var i = 0; i < envs.length; i++) {
            var env = envs[i];
            var confs = env.Configs;
            items.push(
                <div key={env.ID} className="panel">
                    <div className="panel-title">{env.Name}</div>
                    <div className="list">
                    {confs.map((obj) => 
                        <div key={obj.Name} className="list-item">
                            <div className="list-item-title">{obj.Name}: </div>
                            <div className="list-item-desc">{obj.Value}</div>
                        </div>
                    )}
                    </div>
                </div>
            );
        }

        return (
            <div> 
                <table className="header">
                    <tbody>
                        <tr>
                            <td className="header-label">Environments</td>
                            <td className="header-body">
                                <button onClick={this.onEnvCreate} title="Create environment" className="btn-control">
                                    <Add />
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div>{items}</div>
            </div>
        );
    }
}
  
export default Environments;
  