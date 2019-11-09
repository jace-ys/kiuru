import React from "react";

import { AppBar, Toolbar, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  root: {
    "& .MuiToolbar-root": {
      justifyContent: "center",
      "& .MuiTypography-h5": {
        lineHeight: "60px",
        fontWeight: 700
      }
    }
  }
});

interface Props {
  title: string;
}

const TopNav: React.FC<Props> = props => {
  const classes = useStyles();

  return (
    <AppBar position="static" className={classes.root}>
      <Toolbar>
        <Typography variant="h5">{props.title}</Typography>
      </Toolbar>
    </AppBar>
  );
};

export default TopNav;
