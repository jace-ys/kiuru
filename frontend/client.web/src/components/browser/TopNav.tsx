import React from "react";
import { Link } from "react-router-dom";

import { AppBar, Toolbar } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  root: {
    "& .MuiToolbar-root": {
      "& a": {
        height: "80px",
        "& img": {
          width: "40px",
          margin: "20px 10px"
        }
      }
    }
  }
});

const TopNav: React.FC = () => {
  const classes = useStyles();

  return (
    <div className="TopNav">
      <AppBar position="static" color="secondary" className={classes.root}>
        <Toolbar>
          <Link to="/discover">
            <img src={"./assets/icon.png"} alt="Kru" />
          </Link>
        </Toolbar>
      </AppBar>
    </div>
  );
};

export default TopNav;
