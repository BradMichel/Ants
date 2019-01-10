import * as Autosuggest from 'react-autosuggest';
import * as React from "react"
import {FlatButton} from "material-ui"

import {User,Users} from "../class"

interface Create{
    name:string
    ID:string
}

type Suggestion = User|Create
 
interface Props{
    users:Users
    userSelectedID:string
    createUser: (user:User)=>void
    selectUser: (user:User)=>void
}
interface State{
    value:string,
    suggestions:Array<Suggestion>
}

export class MenuUsers extends React.Component<Props,State> {


    constructor() {
        super();
        this.state = {
        value: '',
        suggestions: []
        };
    }

    onChange = (event:any, { newValue }:any) => {
      this.setState({
        value: newValue
      });      
      if(newValue == ""){
        this.props.selectUser(undefined)
      }
    };

    getSuggestions = (value:string) => {
      const {users} = this.props
    
      let usersFound = users.search(value)
      let suggestions:Array<User|Create> = new Array()
      suggestions = suggestions.concat(usersFound)
      suggestions.push({name: value, ID:"-1"})

      return suggestions
    };
    

    getSuggestionValue = (suggestion:Suggestion) => suggestion.name;
    
    renderSuggestion = (suggestion:Suggestion) => {
      if(suggestion.ID === "-1"){
        return(
          <div>
            {suggestion.name}<FlatButton label="Crear Usuario" />
          </div>
        )
      }else{
        return(
          <div>
            {suggestion.name}
          </div>
        )
      }

    };


    onSuggestionsFetchRequested = ({ value }:any) => {
      this.setState({
        suggestions: this.getSuggestions(value)
      });
    };

    onSuggestionsClearRequested = () => {
      this.setState({
        suggestions: []
      });
    };

    onSuggestionSelected = (event:any,{suggestion,suggestionValue,suggestionIndex, method}:{suggestion:Suggestion,suggestionValue:string,suggestionIndex:number,method:string}) => {
      const {createUser,selectUser} = this.props
      if(suggestion.ID === "-1" && method === "click"){
          let user = new User(suggestion.name,"",[],[])
          user.create(createUser)
      }else if(suggestion.ID !== "-1"){
          selectUser(suggestion as User)
      }
      this.setState({value:suggestionValue})
    }

    render() {
      const { value, suggestions } = this.state;

      const inputProps = {
        placeholder: 'Escriba el nombre del usuario',
        value,
        onChange: this.onChange
      };

      return (
        <Autosuggest
          suggestions={suggestions}
          onSuggestionsFetchRequested={this.onSuggestionsFetchRequested}
          onSuggestionsClearRequested={this.onSuggestionsClearRequested}
          getSuggestionValue={this.getSuggestionValue}
          renderSuggestion={this.renderSuggestion}
          inputProps={inputProps}
          onSuggestionSelected={this.onSuggestionSelected}
        />
      );
    }
}
