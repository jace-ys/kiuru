import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import { Hidden } from "@material-ui/core";
import { ThemeProvider, createMuiTheme } from "@material-ui/core/styles";

import BrowserBase from "./components/browser/BrowserBase";
import AppBase from "./components/app/AppBase";
import NotFound from "./pages/NotFound";

import "./App.css";

const theme = createMuiTheme({
  palette: {
    primary: {
      main: "#7a81d9"
    },
    secondary: {
      main: "#ffffff"
    }
  },
  typography: {
    fontFamily: "Futura, Trebuchet MS, Arial, sans-serif"
  }
});

const App: React.FC = () => {
  return (
    <div className="App">
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <Hidden xsDown>
            <BrowserBase />
          </Hidden>
          <Hidden smUp>
            <AppBase />
          </Hidden>
          <Switch>
            <Route path="/discover"></Route>
            <Route path="/connect"></Route>
            <Route path="/notifications"></Route>
            <Route path="/profile"></Route>
            <Route path="/" exact></Route>
            <Route path="*" component={NotFound} />
          </Switch>
        </BrowserRouter>
      </ThemeProvider>
    </div>
  );
};

export default App;
