import React from "react";
import { Link } from "react-router-dom";

import { Button, Grid, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import Base from "../components/Base";
import AppDisplay from "../components/app/Display";
import BrowserDisplay from "../components/browser/Display";

const useStyles = makeStyles(theme => ({
  app: {
    padding: theme.spacing(4),
    display: "flex",
    flexDirection: "column",
    textAlign: "center",
    "& .MuiGrid-root": {
      margin: theme.spacing(4, 0),
      height: "60%"
    },
    "& .MuiButton-root": {
      margin: theme.spacing(1, 0),
      padding: theme.spacing(2, 0)
    }
  },
  browser: {
    padding: theme.spacing(4),
    display: "flex",
    flexDirection: "column"
  }
}));

const AppHome: React.FC = () => {
  const classes = useStyles();

  return (
    <Grid container className={classes.app}>
      <Grid>
        <Typography variant="h4" color="primary">
          Kru Travel
        </Typography>
        <Typography variant="subtitle1">Find your travel crew.</Typography>
      </Grid>
      <Button
        component={Link}
        to="/signup"
        variant="contained"
        color="primary"
        size="large"
        fullWidth
      >
        Sign up
      </Button>
      <Button
        component={Link}
        to="/login"
        variant="contained"
        color="primary"
        size="large"
        fullWidth
      >
        Log in
      </Button>
    </Grid>
  );
};

const BrowserHome: React.FC = () => {
  const classes = useStyles();

  return (
    <Grid container className={classes.browser}>
      <Typography component="h1" variant="h4">
        BrowserHome
      </Typography>
    </Grid>
  );
};

const Home: React.FC = () => {
  return (
    <Base browser>
      <BrowserDisplay>
        <BrowserHome />
      </BrowserDisplay>
      <AppDisplay>
        <AppHome />
      </AppDisplay>
    </Base>
  );
};

export default Home;
