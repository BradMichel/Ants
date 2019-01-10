import * as React from 'react'
import * as Proptypes from 'prop-types'
import Title from 'react-title-component'
import { AppBar, IconButton, FlatButton } from 'material-ui'
import { spacing, getMuiTheme } from 'material-ui/styles'
//import withWidth, {MEDIUM, LARGE} from 'material-ui/utils/withWidth';

import {Ants,AppNavDrawer} from "./"
//import FullWidthSection from './FullWidhSection';

interface Props {
}
interface Context {

}

class Master extends React.Component<Props,Context>{

    static propTypes = {
    }

    static contextTypes = {
        // router: Proptypes.object.isRequired
    }
    state = {
        navDrawerOpen: false,
        coverPageVisibility: true
    }

    handleTouchTapLeftIconButton = () => {
        this.setState({
            navDrawerOpen: !this.state.navDrawerOpen,
        })
    }

    handleChangeRequestNavDrawer = (open:boolean) => {
        this.setState({
            navDrawerOpen: open,
        })
    }

    handleChangeList = (event:any, value:any) =>{
        this.context.router.push(value)
        this.setState({
            navDrawerOpen: false,
        });
    }

    _handleCoverPageVisibility = () => {
        const { coverPageVisibility } = this.state
        this.setState({coverPageVisibility:!coverPageVisibility})

    }

    render(){
        const {  children } = this.props
        let { navDrawerOpen, coverPageVisibility } = this.state
        const router = this.context.router
        let docked:boolean = false;
        let showMenuIconButton:boolean = false; 
        //const styles = getStyles();
        const coverPage = (coverPageVisibility) ? <div style={{display:"inline-block",width:"87%"}}>
                <br /><br /> <br />
                <img style={{}} src="./file/coverPage.png" />
                <FlatButton label="Start" fullWidth={true} onClick={this._handleCoverPageVisibility} />
            </div>:
            <Ants />;
        
        return(
            <div>
                <Title render="INVERTEK S.A. La nueva fuerza" />
                    <AppBar 
                        onLeftIconButtonTouchTap={this.handleTouchTapLeftIconButton}
                        title={"INVERTEK S.A."}
                        style={{position:"fixed"}}
                        zDepth={0}
                        />
                    {coverPage}
                    <AppNavDrawer 
                        docked={docked}
                        onRequestChangeNavDrawer={this.handleChangeRequestNavDrawer}
                        onChangeList={this.handleChangeList}
                        open={navDrawerOpen}
                        />

            </div>
        )

    }

}

export default Master;