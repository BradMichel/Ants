import getMuiTheme from 'material-ui/styles/getMuiTheme'
import {fade} from 'material-ui/utils/colorManipulator';
import {
  deepOrange500, deepOrange50, cyan700,
  pinkA200,
  grey100, grey300, grey400, grey500,
  white, darkBlack, fullBlack,
  redA700
} from 'material-ui/styles/colors';

let getTheme = () =>{

const muiTheme = getMuiTheme({
    fontFamily: 'Roboto, sans-serif',
    palette: {
    primary1Color: deepOrange500,
    primary2Color: cyan700,
    primary3Color: grey400,
    accent1Color: pinkA200,
    accent2Color: grey100,
    accent3Color: grey500,
    textColor: darkBlack,
    alternateTextColor: '#4b4f56',
    canvasColor: white,
    borderColor: grey300,
    disabledColor: fade(darkBlack, 0.3),
    pickerHeaderColor: deepOrange500,
    clockCircleColor: fade(darkBlack, 0.07),
    shadowColor: fullBlack,
    },
    appBar:{
        height: 44,
        color: '#262627',
        textColor:white
    },
    checkbox:{
        checkedColor:"#FFC94D"
    },
    raisedButton:{
        primaryColor: redA700,
        textColor: white,
        primaryTextColor: white,
        secondaryColor: white
    }
})

return muiTheme;
}

export default getTheme;