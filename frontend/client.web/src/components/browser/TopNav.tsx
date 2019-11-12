import React from "react";
import { Link } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";

import { AppBar, Toolbar, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import { RootState } from "../../store";
import { logout } from "../../store/ducks/auth";

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1
  },
  icon: {
    flexGrow: 1,
    height: "80px",
    "& img": {
      height: "40px",
      margin: "20px 10px"
    }
  },
  menu: {
    "& div": {
      display: "flex"
    },
    "& li": {
      display: "inline",
      margin: "0 1rem",
      "& .MuiTypography-root": {
        fontWeight: 700
      },
      "& a": {
        color: theme.palette.primary.main,
        textDecoration: "none"
      },
      "&:hover": {
        textDecoration: "underline"
      }
    }
  }
}));

const TopNav: React.FC = () => {
  const classes = useStyles();

  return (
    <AppBar position="static" color="secondary" className={classes.root}>
      <Toolbar>
        <TopNavIcon />
        <TopNavMenu />
      </Toolbar>
    </AppBar>
  );
};

export default TopNav;

const TopNavIcon: React.FC = () => {
  const classes = useStyles();

  return (
    <Link to="/discover" className={classes.icon}>
      <img src={"./assets/icon.png"} alt="Kru" />
    </Link>
  );
};

const TopNavMenu: React.FC = () => {
  const classes = useStyles();
  const loggedIn = useSelector<RootState, boolean>(
    state => state.auth.loggedIn
  );
  const dispatch = useDispatch();

  return (
    <ul className={classes.menu}>
      {!loggedIn && (
        <div>
          <li>
            <Link to="/signup">
              <Typography variant="h6">Signup</Typography>
            </Link>
          </li>
          <li>
            <Link to="/login">
              <Typography variant="h6">Login</Typography>
            </Link>
          </li>
        </div>
      )}
      {loggedIn && (
        <div>
          <li>
            <Link to="/" onClick={() => dispatch(logout())}>
              <Typography variant="h6">Logout</Typography>
            </Link>
          </li>
        </div>
      )}
    </ul>
  );
};
