import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import { AppBar, Toolbar, IconButton, Typography, MenuItem, Popover } from '@material-ui/core';
import AccountCircle from '@material-ui/icons/AccountCircle';

const styles = theme => ({
  appbar: {
    boxShadow: 'none',
  },
  appbarTitle: {
    flexGrow: 1,
  },
  anchor: {
    color: theme.palette.primary.main,
    transition: '0.7s ease',
    borderTop: '1px solid transparent',
    '&:visited': {
      color: theme.palette.primary.main,
    },
    '&:hover': {
      textShadow: '1px 1px ' + theme.palette.primary.light,
    }
  }
});

class GlobalBar extends Component {

  constructor(props) {
    super(props);
    
    this.state = {
      anchorEl: null,
    };
  }


  handleMenu = event => {
    this.setState({ anchorEl: event.currentTarget });
  };

  handleMenuClose = () => {
    this.setState({ anchorEl: null });
  };

  handleLogout = (event) => {
    this.setState({ anchorEl: null });
    this.props.onLogout(event);
  }

  handleLogin = (event) => {
    this.setState({ anchorEl: null });
    this.props.onLogin(event);
  }

  render() {
    // const profiles = this.state.profiles;
    const { anchorEl } = this.state;
    const { classes, authd } = this.props;
    const menuOpen = Boolean(anchorEl);

    return (
        <AppBar position="static" className={classes.appbar} color="inherit">
          <Toolbar>
            <Typography variant="h5">
              <Link className={classes.anchor} to={'/'}>thrap</Link>
            </Typography>
            <Typography style={{padding:'10px'}}> | </Typography>
            <Typography variant="h5">
              <Link className={classes.anchor} to="/projects">projects</Link>
            </Typography>
            <Typography style={{padding:'10px'}}> | </Typography>
            <Typography variant="h5">
              <Link className={classes.anchor} to={'/docs'}>docs</Link>
            </Typography>
            <Typography className={classes.appbarTitle}></Typography>
            <IconButton color="inherit"
                aria-owns={menuOpen ? 'menu-globalbar' : undefined}
                aria-haspopup="true"
                onClick={this.handleMenu}>
              <AccountCircle />
            </IconButton>
            <Popover
              id="menu-globalbar"
              anchorEl={anchorEl}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              open={menuOpen}
              onClose={this.handleMenuClose}
            >
            {authd ? 
              <MenuItem onClick={this.handleLogout}>
                Logout
              </MenuItem>
              :
              <MenuItem onClick={this.handleLogin}>
                Login
              </MenuItem>
            }
              
            </Popover>
          </Toolbar>
        </AppBar>
    );
  }
}

export default withStyles(styles)(GlobalBar);
