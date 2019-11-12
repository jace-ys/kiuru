import React from "react";
import { useDispatch } from "react-redux";

import { Button, Grid, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import { login } from "../../store/ducks/auth";

const useStyles = makeStyles({
  root: {
    height: "100vh"
  }
});

const Auth: React.FC = () => {
  const classes = useStyles();
  const dispatch = useDispatch();

  return (
    <Grid
      container
      justify="center"
      alignItems="center"
      className={classes.root}
    >
      <Typography>Auth</Typography>
      <Button onClick={() => dispatch(login())}>Login</Button>
    </Grid>
  );
};

export default Auth;
