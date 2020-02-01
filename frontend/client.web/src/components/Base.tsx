import React from "react";
import { Redirect } from "react-router-dom";
import { useSelector } from "react-redux";

import { CssBaseline, Grid } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import AppDisplay from "./app/Display";
import AppBottomNav from "./app/BottomNav";
import AppTopNav from "./app/TopNav";
import BrowserDisplay from "./browser/Display";
import BrowserTopNav from "./browser/TopNav";

import { RootState } from "../store";

const useStyles = makeStyles(theme => ({
  root: {
    width: "100vw",
    height: "100vh",
    [theme.breakpoints.down("xs")]: {
      paddingTop: "56px"
    },
    [theme.breakpoints.up("sm")]: {
      paddingTop: "75px"
    }
  }
}));

interface Props {
  app?: boolean;
  browser?: boolean;
  authenticated?: boolean;
}

const Base: React.FC<Props> = props => {
  const classes = useStyles(props);
  const loggedIn = useSelector<RootState, boolean>(
    state => state.auth.loggedIn
  );

  return (
    <div>
      <CssBaseline />
      {props.app && (
        <AppDisplay>
          {!loggedIn && <Redirect to="/"></Redirect>}
          <AppTopNav />
          <AppBottomNav />
        </AppDisplay>
      )}
      {props.browser && (
        <BrowserDisplay>
          {props.authenticated && !loggedIn && (
            <Redirect to="/login"></Redirect>
          )}
          <BrowserTopNav />
        </BrowserDisplay>
      )}
      <Grid container className={classes.root}>
        {props.children}
      </Grid>
    </div>
  );
};

export default Base;
