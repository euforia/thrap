import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {Typography, Checkbox, FormControlLabel, Divider} from '@material-ui/core';
import Modal from '@material-ui/core/Modal';
import { Button } from '@material-ui/core';


const styles = theme => ({
  paper: {
    position: 'absolute',
    width: theme.spacing.unit * 50,
    backgroundColor: theme.palette.background.paper,
    boxShadow: theme.shadows[5],
    padding: theme.spacing.unit * 4,
    top: '50%',
    left: '50%',
    transform: `translate(-50%, -50%)`,
    borderRadius: '4px',
  },
  body: {
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
  },
  footer: {
    textAlign: 'right',
    paddingTop: theme.spacing.unit * 3,
  },
  button: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
  purge: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
  },
  divider: {
    marginTop: theme.spacing.unit * 4,
  }
});

class DeployStop extends Component {

    state = {
        purge: false,
    }

    handleCheck = (event) => {
        var p = event.target.checked;
        this.setState({purge: p});
    }
    handleCancel = (event) => {
        this.setState({purge:false});
        this.props.onCancel(event);
    }

  render() {
    const { classes } = this.props;
    const {purge} = this.state;

    return (
      <div>
        <Modal
          open={this.props.open}
          onClose={this.props.onCancel}
        >
          <div className={classes.paper}>
            <Typography variant="h5">
              Are you sure you want to stop <b>{this.props.name}</b>?
            </Typography>
            <Divider className={classes.divider}/>
            <FormControlLabel className={classes.purge}
                control={<Checkbox
                    color="secondary"
                    checked={purge}
                    onChange={this.handleCheck}
                />}
                label="Purge (remove record from orchestrator)"
            />
            <Divider/>
            <div className={classes.footer}>
                <Button color="default" className={classes.button} 
                  onClick={this.handleCancel}>Cancel</Button>
                <Button color="secondary" className={classes.button}
                  onClick={event => this.props.onStop(purge, event)}>Stop</Button>
            </div>
          </div>
        </Modal>
      </div>
    );
  }
}

DeployStop.propTypes = {
  classes: PropTypes.object.isRequired,
  onCancel: PropTypes.func.isRequired,
  onStop: PropTypes.func.isRequired,
  open: PropTypes.bool.isRequired,
};

export default withStyles(styles)(DeployStop);
