import React from "react";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";

import { ThemeProvider, createMuiTheme } from "@material-ui/core/styles";

import Home from "./pages/Home";
import Signup from "./pages/Signup";
import Login from "./pages/Login";
import Discover from "./pages/Discover";
import NotFound from "./pages/NotFound";

import "./App.css";

const theme = createMuiTheme({
  palette: {
    primary: {
      main: "#7a81d9"
    },
    secondary: {
      main: "#ffffff",
      contrastText: "#484848"
    },
    background: {
      default: "#ffffff"
    },
    text: {
      primary: "#484848",
      secondary: "#484848"
    }
  },
  typography: {
    fontFamily: "Futura, Tahoma, Arial, sans-serif",
    h4: {
      fontFamily: "Futura-Bold, Tahoma, Arial, sans-serif",
      fontWeight: 700
    },
    h6: {
      fontFamily: "Futura-Bold, Tahoma, Arial, sans-serif"
    },
    button: {
      fontFamily: "Futura-Bold, Tahoma, Arial, sans-serif",
      textTransform: "none"
    }
  }
});

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <Route
          component={() => {
            window.scrollTo(0, 0);
            return null;
          }}
        />
        <Switch>
          <Route path="/" exact component={Home}></Route>
          <Route path="/signup" exact component={Signup}></Route>
          <Route path="/login" exact component={Login}></Route>
          <Route
            path="/logout"
            exact
            component={() => <Redirect to="/"></Redirect>}
          ></Route>
          <Route path="/discover" exact component={Discover}></Route>
          <Route path="/connect" exact></Route>
          <Route path="/notifications" exact></Route>
          <Route path="/profile" exact></Route>
          <Route path="*" component={NotFound} />
        </Switch>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App;
