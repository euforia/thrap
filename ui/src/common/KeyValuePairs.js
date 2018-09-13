import React, { Component } from 'react';
import { Clear, Add } from '@material-ui/icons';

const styles = ({
    container: {
        fontWeight: 100,
    },
    icon: {
        height: 19,
        width: 19,
    },
    title: {
        padding: '10px 5px',
        // borderBottom: '1px solid #4b5f6d',
        display: 'inline-block',
    },
})

class KeyValuePairs extends Component {
    constructor(props) {
        super(props);

        this.addKVPair = this.addKVPair.bind(this);
        this.removeKVPair = this.removeKVPair.bind(this);
        this.onKeyChange = this.onKeyChange.bind(this);
        this.onValueChange = this.onValueChange.bind(this);

        this.state = {
            pairs: props.pairs,
        }
    }

    addKVPair() {
        var pairs = this.state.pairs;
        pairs.push({name:'',value:''});
        this.setState(pairs);
    }

    removeKVPair(event) {
        var pairs = this.state.pairs;
        var name = event.currentTarget.name;
        if (pairs[name].required) {
            return 
        }

        pairs.splice(name, 1);
        this.setState({pairs: pairs});

        this.props.onKVChange(pairs);
    }

    onKeyChange(event) {
        var i = event.currentTarget.name;
        var pairs = this.state.pairs;
        pairs[i].name = event.currentTarget.value;
        pairs[i].keyClass = pairs[i].name === '' ? 'invalid-input' : '';
        this.setState({pairs:pairs});
    }

    onValueChange(event) {
        var i = event.currentTarget.name;
        var pairs = this.state.pairs;
        pairs[i].value = event.currentTarget.value;
        pairs[i].valueClass = pairs[i].value === '' ? 'invalid-input' : '';
        this.setState({pairs:pairs});

        this.props.onKVChange(pairs);
    }

    render() {
        var pairs = this.state.pairs;

        return (
            <div style={styles.container}>
                <table>
                    <thead>
                        <tr>
                            <td className="theme-color" style={{padding: '20px 0'}}>
                                <div style={styles.title} className="vars-header">{this.props.title}</div>
                            </td>
                            <td style={{textAlign: "center"}}>
                                <button title="Add" className='btn-trans' onClick={this.addKVPair}>
                                    <span className="btn-datastore-icon">
                                        <Add style={styles.icon}/>
                                    </span>
                                </button>
                            </td>
                        </tr>
                    </thead>
                    <tbody>
                        {pairs.map((obj, keyIndex) => 
                        <tr key={keyIndex}>
                            <td className="theme-color">
                                <input type="text" className={obj.keyClass} name={keyIndex} value={obj.name} onChange={this.onKeyChange} disabled={obj.required ? 'disabled': ''} /> 
                                = 
                                <input type="text" name={keyIndex} value={obj.value} onChange={this.onValueChange} className={obj.valueClass}/>
                            </td>
                            <td style={{textAlign: "center"}}>
                                <button title="Remove" className={obj.required ? 'hide' : 'btn-trans' } onClick={this.removeKVPair} name={keyIndex}>
                                    <span className="btn-datastore-icon">
                                        <Clear style={styles.icon}/>
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

export default KeyValuePairs;