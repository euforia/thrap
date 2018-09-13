import React, { Component } from 'react';
import { ExpandMore, Add, Clear } from '@material-ui/icons';
import './DatastoreSelector.css'

class DatastoreSelector extends Component {
   
    constructor(props) {
        super(props);

        this.addDatastore = this.addDatastore.bind(this);
        this.removeDatastore = this.removeDatastore.bind(this);
        this.onSelectChange = this.onSelectChange.bind(this);

        this.state = {
            datastores: [],
        }
    }

    onSelectChange(event) {
        var ds = this.state.datastores;
        ds[event.target.name] = event.target.value; 
        this.setState({datastores: ds});
    };

    addDatastore() {
        var ds = this.state.datastores;
        ds.push("");
        this.setState({datastores: ds});
    }

    removeDatastore(event) {
        var ds = this.state.datastores;
        ds.splice(event.currentTarget.name, 1);
        this.setState({datastores: ds});
    }

    render() {
        var datastores = this.state.datastores;

        return (
            <table id="datastores">
                <tbody>
                    <tr>
                        <td className="section-label datastore-label-cell">Datastore</td>
                        <td className="datastore-ctrl-cell">
                            <button className="btn-datastore" onClick={this.addDatastore}>
                                <span className="btn-datastore-icon">
                                    <Add style={{height: 20, width: 20}}/>
                                </span>
                            </button>
                        </td>
                    </tr>
                    {datastores.map((obj, keyIndex) =>
                        <tr key={keyIndex}>
                            <td>
                                <div className="select-container">
                                    <select value={obj} onChange={this.onSelectChange} name={keyIndex}>
                                        <option value="elasticsearch">ElasticSearch</option>
                                        <option value="postgres">Postgres</option>
                                        <option value="mysql">MySQL</option>
                                        <option value="redis">Redis</option>
                                    </select>
                                    <span className="select-expand">
                                        <ExpandMore/>
                                    </span>
                                </div>
                            </td>
                            <td className="datastore-ctrl-cell">
                                <button className="btn-datastore" onClick={this.removeDatastore} name={keyIndex}>
                                    <span className="btn-datastore-icon">
                                        <Clear style={{height: 19, width: 19}}/>
                                    </span>
                                </button>
                            </td>
                        </tr>
                    )}

                </tbody>
            </table>
        );
    }
}
  
export default DatastoreSelector;