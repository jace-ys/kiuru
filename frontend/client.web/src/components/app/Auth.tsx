import React from "react";

import { Button, Grid, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  root: {
    height: "100vh"
  }
});

interface Props {
  loggedIn: boolean;
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

const Auth: React.FC<Props> = props => {
  const classes = useStyles();
  const { loggedIn, setLoggedIn } = props;

  return (
    <Grid
      container
      justify="center"
      alignItems="center"
      className={classes.root}
    >
      <Typography>Auth</Typography>
      <Button onClick={() => setLoggedIn(!loggedIn)}>Login</Button>
    </Grid>
  );
};

export default Auth;
