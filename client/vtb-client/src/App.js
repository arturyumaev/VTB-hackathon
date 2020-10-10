import React from 'react';
import './App.css';
import { BrowserRouter, Route, Switch } from 'react-router-dom'
import Home from './Pages/Home/Home';
import { Login } from './Pages/Login/Login';
import { Signup } from './Pages/Signup/Signup';
import Navbar from './Components/Navbar';
import SuspectActivity from './Pages/SuspectActivity/SuspectActivity';
import Portal from './Pages/Portal/Portal';

function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <div className="container">
        <Switch>
          <Route path={'/'} exact component={Portal} />
          <Route path={'/login'} component={Login} />
          <Route path={'/signup'} component={Signup} />
          <Route path={'/suspect_activity'} component={SuspectActivity} />
          <Route path={'/portal'} component={Portal} />
        </Switch>
      </div>
    </BrowserRouter>
  );
}

export default App;
