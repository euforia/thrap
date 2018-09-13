import React, { Component } from 'react';
import { ExpandMore } from '@material-ui/icons';

class BuildConfig extends Component {

    constructor(props) {
        super(props);

        this.state = {
            configs: props.configs,
        }

    }

    render() {
        return (
            <table>
                <tbody>
                    {this.state.configs.map((conf) => 
                    <tr key={conf.id}>
                        <td>{conf.id}</td>
                        <td>
                            <div>Language</div>
                            <div className="create-proj-input select-container">
                                <select value={conf.language}>
                                    <option value="none">None</option>
                                    <option value="shell">Shell</option>
                                    <option value="go">Golang</option>
                                    <option value="python">Python</option>
                                    <option value="java">Java</option>
                                    <option value="ruby">Ruby</option>
                                    <option value="csharp">C#</option>
                                </select>
                                <span className="select-expand">
                                    <ExpandMore/>
                                </span>
                            </div>
                        </td>
                        <td>
                            <div>Dockerfile</div>
                            <div><input type="text" value={conf.build.dockerfile}/></div>
                        </td>
                        <td>
                            <div>Context</div>
                            <div>
                                <input type="text" value={conf.build.context}/>
                            </div>
                        </td>
                    </tr>
                    )}
                </tbody>
            </table>
        );
    }
}

export default BuildConfig;