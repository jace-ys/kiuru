import React from "react";

import { AppBar, Toolbar, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  root: {
    "& .MuiToolbar-root": {
      justifyContent: "center"
    }
  }
});

const TopNav: React.FC = () => {
  const classes = useStyles();

  return (
    <AppBar position="fixed" className={classes.root}>
      <Toolbar>
        <Typography variant="h6">Kru Travel</Typography>
      </Toolbar>
    </AppBar>
  );
};

export default TopNav;
