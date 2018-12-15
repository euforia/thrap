import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
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
  },
  body: {
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
  },
  footer: {
    textAlign: 'right'
  },
  button: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  }
});

class ConfirmDelete extends Component {

  render() {
    const { classes } = this.props;

    return (
      <div>
        <Modal
          open={this.props.open}
          onClose={this.props.onCancel}
        >
          <div className={classes.paper}>
            <Typography variant="h5">
              Delete {this.props.entityType}
            </Typography>
            <Typography className={classes.body}>
                Are you sure you want to delete <b>{this.props.entity}</b>?
            </Typography>
            <div className={classes.footer}>
                <Button color="default" className={classes.button} 
                  onClick={this.props.onCancel}>Cancel</Button>
                <Button color="secondary" className={classes.button} variant="contained" 
                  onClick={this.props.onDelete}>Delete</Button>
            </div>
          </div>
        </Modal>
      </div>
    );
  }
}

ConfirmDelete.propTypes = {
  classes: PropTypes.object.isRequired,
  onCancel: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default withStyles(styles)(ConfirmDelete);
