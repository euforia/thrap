import React, { Component } from 'react';
import { Grid } from "@material-ui/core";
import { Clear, ImportExportOutlined } from '@material-ui/icons';

const styles = ({
    icon: {
        fontSize: 44,
    },
    iconContainer: {
        padding: '26px',
        textAlign: 'center',
    },
});

class SpecControls extends Component {

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
                    <div style={styles.iconContainer} title="Import/Export">
                        <button className="btn-trans" onClick={this.props.onImportExport}>
                            <ImportExportOutlined style={styles.icon} title="Specification"/>
                        </button>
                    </div>
                </Grid>
            </Grid>
        );
    }
}

export default SpecControls;