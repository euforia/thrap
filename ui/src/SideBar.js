import React, { Component } from 'react';
import { Apps, PersonPin, Code, MoreHoriz } from '@material-ui/icons';
import { withStyles } from '@material-ui/core/styles';
import { IconButton } from '@material-ui/core';

import './SideBar.css';

import thrap from './thrap';
import { Link } from 'react-router-dom';

const styles = theme => ({
  button: {
    margin: theme.spacing.unit,
  },
});

class SideBar extends Component {
    constructor(props) {
      super(props);

      this.signOut = this.signOut.bind(this);
    }

    signOut() {
      thrap.deauthenticate();
      this.props.onSignOut();
    }

    render() {
        const { classes } = this.props;
        return (
            <div id="sidebar">
              <table style={{height: '300px'}}>
                  <tbody>
                    <tr>
                      <td style={{textAlign: 'center'}}>
                        <div >
                          <PersonPin style={{height: 150, width:150}} />
                        </div>
                        <div>
                          <IconButton className={classes.button} color="primary" 
                            aria-label="Log Out" title="Log Out"
                            onClick={this.signOut}>
                            <MoreHoriz />
                          </IconButton>
                        </div>
                      </td>
                    </tr>
                  </tbody>
              </table>
              <div title="Projects" name="projects">
                <Link to="/projects">
                  <div className="sidebar-item">
                    <Apps />
                    <span className="sidebar-item-label">Projects</span>
                  </div>
                </Link>
              </div>
              <div title="Documentation" name="docs">
                <Link to="/docs">
                  <div className="sidebar-item">
                    <Code />
                    <span className="sidebar-item-label">Documentation</span>
                  </div>
                </Link>
              </div>
            </div>
        );
    }
  }
  
  export default withStyles(styles)(SideBar);
  // export default SideBar;
  