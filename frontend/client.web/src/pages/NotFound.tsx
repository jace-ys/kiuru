import React from "react";

import { Grid, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import Base from "../components/Base";

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
    display: "flex",
    flexDirection: "column"
  }
}));

const NotFound: React.FC = () => {
  const classes = useStyles();

  return (
    <Base app browser>
      <Grid container className={classes.root}>
        <Typography component="h1" variant="h4">
          Not Found
        </Typography>
      </Grid>
    </Base>
  );
};

export default NotFound;
