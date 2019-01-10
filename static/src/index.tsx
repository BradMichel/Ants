import * as React from 'react';
import {render} from 'react-dom';
import Master from './components/Master';
import {createHashHistory} from 'history';
import injectTapEventPlugin = require("react-tap-event-plugin");
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getTheme from './styles/muiTheme';

// window.React = React;
//window.Perf = require('react-addons-perf');

// require('./styles/main')
injectTapEventPlugin();

render(
  <MuiThemeProvider muiTheme={getTheme()}>
   <Master />
  </MuiThemeProvider>,
  document.getElementById('app')
);
