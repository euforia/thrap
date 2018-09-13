import React, { Component } from 'react';
import { SettingsSharp, Build, Clear } from '@material-ui/icons';
import './Builds.css'
import BuildConfig from './BuildConfig';

class Builds extends Component {

    constructor(props) {
        super(props);

        this.onConfigure = this.onConfigure.bind(this);
        // this.onCloseConfigure = this.onCloseConfigure.bind(this);

        this.state = {
            configure: false,
            configs: props.configs,
            builds: [
                // {
                //     id: '3',
                //     commit: '8a2dfb233db06a1ef4c844c6b41acd202abc123',
                //     commitMessage: "Third commit message test",
                //     artifacts: ["project/api:v0.1.2-2-abc123", "project/web:v0.3-2-abc123"],
                //     status: "running",
                //     runtime: '0',
                //     log: '',
                // },{
                //     id: '2',
                //     commit: '9a2dfb233db06a1ef4c844c6b41acd202cf3c993',
                //     commitMessage: "Second commit message test",
                //     artifacts: ["project/api:v0.1.2", "project/web:v0.3"],
                //     status: "succeeded",
                //     runtime: '3m2s',
                //     log: '',
                // },{
                //     id: '1',
                //     commit: '1ef4c844c6b419a2dfb233db06aacd202cf3c993',
                //     commitMessage: "Initial commit",
                //     artifacts: [],
                //     status: "failed",
                //     runtime: '32s',
                //     log: '',
                // }
            ],
        }
    }

    onConfigure() {
        var c = !this.state.configure;
        this.setState({configure: c});
    }
    // onCloseConfigure() {
    //     this.setState({configure:false});
    // }


    render() {
        var configs = [];
        if (this.state.configure) {
            configs.push(
                <tr>
                    <td colSpan="4">
                        <BuildConfig configs={this.state.configs} />
                    </td>
                </tr>
            );
        }

        return (
            <div>
                <table id="builds-table">
                    <thead>
                        <tr>
                            <td className="build-id"></td>
                            <td></td>
                            <td colSpan="2" style={{textAlign:"right"}}>
                                <button title="Configure" style={{margin: "0 10px"}} className="btn-control" onClick={this.onConfigure}>
                                    <SettingsSharp />
                                </button>
                                <button title="Build" style={{margin: "0 10px"}} className="btn-control"><Build /></button>
                            </td>
                        </tr>
                        {configs}
                        <tr>
                            <td className="build-id"><div>#</div></td>
                            <td></td>
                            <td className="build-artifacts"><div>Artifacts</div></td>
                            <td className="build-status"><div>Status</div></td>
                        </tr>
                    </thead>
                    <tbody>
                    {this.state.builds.map((obj) => 
                        <tr key={obj.id}>
                            <td className="build-id">{obj.id}</td>
                            <td className="build-commit">
                                <div>{obj.commitMessage}</div>
                                <div className="subscript">{obj.commit}</div>
                            </td>
                            <td className="build-artifacts">
                                {obj.artifacts.map((art) => 
                                    <div key={art}>{art}</div>
                                )}
                            </td>
                            <td className="build-status">
                                <div className={obj.status === 'failed' ? 'build-status-failed' : ''}>{obj.status}</div>
                                <div className="subscript">{obj.runtime}</div>
                            </td>
                            <td style={{textAlign: "center"}}>
                                <button title="Cancel" className={obj.status === 'running' ? 'btn-datastore' : 'btn-datastore btn-hide'}>
                                    <span className="btn-datastore-icon">
                                        <Clear style={{height: 19, width: 19}}/>
                                    </span>
                                </button>
                            </td>
                        </tr>
                    )}
                    </tbody>
                </table>
            </div>
        );
    }
  }
  
  export default Builds;
  