import * as React from 'react';
import {
    Route, Redirect, IndexRoute
} from 'react-router';

import Master from './components/Master';
import Ants from './components/Ants'


const AppRoutes = (
    <Route path="/" component={Ants}>
    <IndexRoute component={Ants} />
    <Route path="ants" component={Ants} />
    </Route > 
)

export default AppRoutes;