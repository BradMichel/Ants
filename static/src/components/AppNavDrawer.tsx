import * as React from 'react'
import * as PropTypes from 'prop-types'
import { Drawer, List, ListItem, makeSelectable, Divider, Subheader, DropDownMenu, MenuItem } from 'material-ui'
import { spacing, typography, zIndex } from 'material-ui/styles'

const SelectableList = makeSelectable(List)

interface Props{
    docked: boolean
    onChangeList: Function
    onRequestChangeNavDrawer: (opening: boolean, reason: string) => void
    open: boolean
}
interface Context{
    // router: Object
}

export class AppNavDrawer extends React.Component<Props,Context>{
    static propTypes = {
        docked: PropTypes.bool.isRequired,
        onChangeList: PropTypes.func.isRequired,
        onRequestChangeNavDrawer: PropTypes.func.isRequired,
        open: PropTypes.bool.isRequired,
    }
    static contextTypes:Context={
        // router: PropTypes.object.isRequired
    }
    handleTouchTapHeader = () =>{
        // this.context.router.push('/')
        this.props.onRequestChangeNavDrawer(false, '')
    }

    render(){
        const {
             docked, onRequestChangeNavDrawer, onChangeList, open 
        } = this.props

        return(
            <Drawer 
                docked={docked}
                open={open}
                onRequestChange={onRequestChangeNavDrawer}
                 >

            </Drawer>
        )
    }
}
