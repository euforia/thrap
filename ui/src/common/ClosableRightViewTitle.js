import React, { Component } from 'react';
import { Grid, IconButton } from "@material-ui/core";
import { Clear } from '@material-ui/icons';

const styles = ({
    // title: {
    //     fontSize: '24px',
    //     fontWeight: 100,
    // },
    subtext: {
        fontSize: '12px',
        opacity: 0.7,
        fontWeight: 100,
        paddingTop: '5px',
    },
    // container: {
    //     padding: '20px 0',
    // }
});

class ClosableRightViewTitle extends Component {
    constructor(props) {
        super(props);

        this.state = {
            title: props.title === undefined ? '' : props.title,
        }
    }

    getSubText() {
        if (this.props.subtext === undefined) {
            return;
        } else if (this.props.subtext === '') {
            return;
        }

        return (
            <div style={styles.subtext}>{this.props.subtext}</div>
        )
    }

    render() {
        return (
            <div className="header-container">
                <Grid container spacing={0} alignItems="center">
                    <Grid item xs={9}>
                        <div className="header-title">{this.state.title}</div> 
                        {this.getSubText()}
                    </Grid>
                    <Grid item xs={3}>
                        <div style={{textAlign: 'right'}}>
                            <IconButton size="large" onClick={this.props.onCloseDialogue} aria-label="Close">
                                <Clear />
                            </IconButton>
                            {/* <button title="Close" className="btn-control" onClick={this.props.onCloseDialogue}>
                                <Clear />
                            </button> */}
                        </div>
                    </Grid>
                </Grid>
            </div>
        );
    }
}

export default ClosableRightViewTitle;