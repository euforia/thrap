import React, { Component } from 'react';
import { Grid } from "@material-ui/core";
import { Clear, Cached, SettingsOutlined } from '@material-ui/icons';

const styles = ({
    icon: {
        fontSize: 44,
    },
    iconContainer: {
        padding: '26px',
        textAlign: 'center',
    },
});

class DeployControls extends Component {
    render() {
        return (
            <Grid container spacing={0} direction="column" alignItems="flex-start">
                <Grid item xs={12}>
                    <div style={styles.iconContainer} title="Close">
                        <button className="btn-trans" onClick={this.props.onClose}>
                            <Clear style={styles.icon} />
                        </button>
                    </div>
                </Grid>
                <Grid item xs={12}>
                    <div style={styles.iconContainer}>
                        <button className="btn-trans" onClick={this.props.onSettings}>
                            <SettingsOutlined style={styles.icon} title="Specification"/>
                        </button>
                    </div>
                </Grid>
                <Grid item xs={12}>
                    <div style={styles.iconContainer} title="Deploy">
                        <button className='btn-trans' onClick={this.props.onDeploy}>
                            <Cached style={styles.icon} />
                        </button>
                    </div>
                </Grid>
            </Grid>
        );
    }
}

export default DeployControls;