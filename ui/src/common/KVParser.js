import React, { Component } from 'react';
import { IconButton, InputAdornment, TextField } from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';


function validateKV(kv) {
    var pair = kv.split('=');
    if (pair.length!==2) {
        return null;
    }
    pair[0] = pair[0].trim();
    pair[1] = pair[1].trim();
    if (pair[0].length===0 || pair[1].length ===0) {
        return null;
    }
    return pair;
}

class KVParser extends Component {
    constructor(props) {
        super(props);

        this.state = {
            keyName: '',
            keyErr: false,
        }
    }

    addKVPair = () => {
        var pair = validateKV(this.state.keyName)
        if (pair === null) {
            this.setState({keyErr:true});
            return;
        }

        this.setState({keyName:'', keyErr:false});
        this.props.onKVAdd(pair[0], pair[1]);
    }

    handleChange = event => {
        var val = event.currentTarget.value;
        var pair = validateKV(val);

        this.setState({keyName: val, keyErr: pair===null});
    }
    
    onKeyPress = (event) => {
        if (event.key === 'Enter') {
            this.addKVPair();
        }
    }

    render() {
        return (              
            <div style={{margin: '15px 0', paddingTop: '10px'}}>
                <TextField
                    helperText="Optionally add runtime variable"
                    fullWidth
                    onKeyPress={this.onKeyPress}
                    variant="outlined"
                    margin="dense"
                    value={this.state.keyName}
                    onChange={this.handleChange}
                    error={this.state.keyErr}
                    placeholder="key=value"
                    InputProps={{
                        endAdornment: 
                            <InputAdornment>
                                <IconButton color="primary" onClick={this.addKVPair}>
                                    <AddIcon fontSize="small"/>
                                </IconButton>
                            </InputAdornment>,
                    }}
                />
            </div>
        );
    }
}

export default (KVParser);