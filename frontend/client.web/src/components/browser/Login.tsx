import React from "react";

import { Button, Grid, Typography } from "@material-ui/core";

interface Props {
  loggedIn: boolean;
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

const Login: React.FC<Props> = props => {
  const { loggedIn, setLoggedIn } = props;

  return (
    <Grid container justify="center" alignItems="center">
      <Typography>Login</Typography>
      <Button onClick={() => setLoggedIn(!loggedIn)}>Login</Button>
    </Grid>
  );
};

export default Login;
