import React from "react";
import { useDispatch } from "react-redux";

import { Button, Grid, Typography } from "@material-ui/core";

import { login } from "../../store/ducks/auth";

const Login: React.FC = () => {
  const dispatch = useDispatch();

  return (
    <Grid container justify="center" alignItems="center">
      <Typography>Login</Typography>
      <Button onClick={() => dispatch(login())}>Login</Button>
    </Grid>
  );
};

export default Login;
