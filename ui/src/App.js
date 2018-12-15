import React, { Component } from 'react';
import { Switch, Route, withRouter, Redirect } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';

import Login from './login/Login';
import NewProject from './project/NewProject';
import Project from './project/Project';
import Projects from './project/Projects';
import APIDocs from './docs/Docs';
import GlobalBar from './common/GlobalBar';

import thrap from './api/thrap';

const defaultRedirect = "/";
const styles = theme => ({
  stageCont: {
    paddingLeft: theme.spacing.unit * 20,
    paddingRight: theme.spacing.unit * 20,
  }
});

class App extends Component {

  constructor(props) {
    super(props);
    
    this.state = {
      profiles: [],
      authd: thrap.isAuthd(),
    };
    
    this.fetchProfiles();
  }

  fetchProfiles() {
    thrap.Profiles().then(profs => {
      this.setState({
        profiles: profs.data,
      });
    });
  }

  onLoginSucceeded = (data, event) => {
    this.setState({authd:true});

    var nextPath = this.props.location.hash;
    if (nextPath === '') nextPath = defaultRedirect;
    else nextPath = nextPath.replace('#', '');

    this.props.history.push(nextPath);
  }

  handleLogout = (event) => {
    thrap.Deauthenticate();
    this.setState({authd:false});
  }

  redirectToLogin = () => {
    var to = '/login#'+this.props.location.pathname;
    this.props.history.push(to);
  }

  render() {
    const { profiles, authd } = this.state;
    const { classes } = this.props;

    return (
      <div>
        <GlobalBar 
          onLogin={this.redirectToLogin}
          onLogout={this.handleLogout} 
          authd={authd}
        />
        <Switch>
          <Route path="/docs" 
            render={(props) => <APIDocs {...props} />} 
          />
          <Route path="/login" 
            render={(props) => <Login {...props} profiles={profiles} onLogin={this.onLoginSucceeded} />}
          />
          <Route path="/projects" exact 
            render={(props) => 
              <div className={classes.stageCont}>
                <Projects {...props} />
              </div>
            } 
          />
          <Route path="/projects/new" exact 
            render={(props) => 
              <div className={classes.stageCont}>
                <NewProject {...props} />
              </div>
            } 
          />
          <Route path="/project/:project" 
            render={(props) => 
              <div className={classes.stageCont}>
                <Project {...props} profiles={profiles} />
              </div>
            } 
          />
          <Route path="/" 
            render={(props) => <Redirect to="/projects" {...props} />} 
          />
        </Switch>
      </div>
    );
  }
}

export default withRouter(withStyles(styles, { withTheme: true })(App));
