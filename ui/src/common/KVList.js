import React, { Component } from 'react';
import { withStyles, IconButton, InputAdornment, TextField } from '@material-ui/core';
import ClearIcon from '@material-ui/icons/Clear';
import KVParser from '../common/KVParser';
// import AddIcon from '@material-ui/icons/Add';


const styles = theme => ({
    container: {
        paddingBottom: 10,
    },
    icon: {
        height: 19,
        width: 19,
    },
})

class KVList extends Component {

    removeKVPair = event => {
        var i = event.currentTarget.name;
        var pairs = this.props.pairs;
        if (pairs[i].required) {
            return 
        }
        
        this.props.onKVRemove(i, event);
    }

    inputAdornment(ro, keyIndex) {
        var pairs = this.props.pairs;
        // var pairs = this.state.pairs;
        if (ro === true || pairs[keyIndex].required === true) return null;
        return (
            <InputAdornment position="start">
                <IconButton 
                    name={keyIndex}
                    onClick={this.removeKVPair}>
                    <ClearIcon fontSize="small"/>
                </IconButton>
            </InputAdornment>
        );
    }

    render() {
        const { classes, readOnly, disabled } = this.props;
        const {pairs} = this.props;

        return (
            <div className={classes.container}>
                {pairs.map((obj, keyIndex) => 
                    <TextField 
                        key={keyIndex}
                        name={obj.name}
                        label={obj.name}
                        value={obj.value}
                        margin="dense"
                        onChange={event => this.props.onKVChange(keyIndex, event)}
                        fullWidth
                        required
                        disabled={disabled}
                        InputProps={{
                            readOnly: readOnly,
                            startAdornment: this.inputAdornment(readOnly, keyIndex),
                        }}
                        error={obj.error}
                    />
                )}
                {readOnly ? <div/> : <KVParser onKVAdd={this.props.onKVAdd} />}
            </div>
        );
    }
}

export default withStyles(styles)(KVList);