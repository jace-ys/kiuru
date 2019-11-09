import React from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";

import { Hidden } from "@material-ui/core";
import { ThemeProvider, createMuiTheme } from "@material-ui/core/styles";

import AppBase from "./components/app/Base";

import BrowserBase from "./components/browser/Base";
import BrowserLogin from "./components/browser/Login";
import BrowserSignup from "./components/browser/Signup";

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
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <Hidden xsDown>
          <BrowserBase />
        </Hidden>
        <Hidden smUp>
          <AppBase />
        </Hidden>
        <Switch>
          <Route path="/" exact></Route>

          <Route path="/signup" component={BrowserSignup}></Route>
          <Route path="/login" component={BrowserLogin}></Route>
          <Route path="/logout" render={() => <Redirect to="/" />} />

          <Route path="/discover"></Route>
          <Route path="/connect"></Route>
          <Route path="/notifications"></Route>
          <Route path="/profile"></Route>
          <Route path="*" component={NotFound} />
        </Switch>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App;
