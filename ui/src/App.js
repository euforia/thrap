import React, { Component } from 'react';
import { ChevronLeft, ChevronRight } from '@material-ui/icons';
import Login from './Login.js';
import SideBar from './SideBar.js';
import Projects from './Projects.js';
// import Envrionments from './Envrionments.js';
// import thrap from './thrap';
import APIDocs from './APIDocs';

import './App.css';
import { Switch, Redirect } from 'react-router-dom';
import Route from 'react-router-dom/Route';

class App extends Component {
  constructor(props) {
    super(props);

    this.onLoginSucceeded = this.onLoginSucceeded.bind(this);
    this.onSignOut = this.onSignOut.bind(this);

    this.toggleSideBar = this.toggleSideBar.bind(this);

    this.state = {
      stageClass: 'hide-sidebar',
      loggedIn: false,
    };

  }

  onLoginSucceeded() {
    console.log('Login succeeded');

    this.setState({
      loggedIn: true,
      stageClass: 'close-sidebar',
    });
  }

  onSignOut() {
    this.setState({
      loggedIn: false,
      stageClass: 'hide-sidebar',
    });
  }

  toggleSideBar() {
    if (this.state.stageClass === '') {
      this.setState({stageClass: 'close-sidebar'});
    } else {
      this.setState({stageClass:''});
    }
  }

  render() {
    if (!this.state.loggedIn) {
      return (          
        <div id="stage" className={this.state.stageClass}>
            <Login onLoginSucceeded={this.onLoginSucceeded}></Login>
        </div>
      );

    }
    
    // Temporary hack. Revisit at some point
    var stageClass = this.state.stageClass;
    if (window.location.pathname.includes("docs")) {
      stageClass += " docs";
    }

    return (
      <div>
        <SideBar onSignOut={this.onSignOut} />

        <div id="stage" className={stageClass}>
          <ChevronLeft onClick={this.toggleSideBar} 
            className={"chevron-left-icon " + (this.state.stageClass === '' ? '' : 'hide')} 
          />
          
          <ChevronRight onClick={this.toggleSideBar} 
            className={"chevron-right-icon " + (this.state.stageClass === 'close-sidebar' ? '' : 'hide')} 
          />
         
          <Switch>
            <Route exact path='/projects' component={Projects} />
            <Route path='/project/:project' component={Projects} />
            <Route path='/docs' component={APIDocs} />
            <Redirect from='/' to='/projects'/>
          </Switch>
        </div>

      </div>
    );
  }
}

export default App;
