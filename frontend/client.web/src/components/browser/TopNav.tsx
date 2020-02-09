import React from "react";
import { Link } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";

import { AppBar, Button, Toolbar, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";

import logo from "../../assets/logo.png";

import { RootState } from "../../store";
import { logout } from "../../store/ducks/auth";

const useStyles = makeStyles(theme => ({
  icon: {
    flexGrow: 1,
    height: "75px",
    "& img": {
      height: "40px",
      margin: "17.5px"
    }
  },
  menu: {
    "& .MuiButton-root": {
      margin: theme.spacing(1)
    }
  }
}));

const TopNav: React.FC = () => {
  const classes = useStyles();
  const loggedIn = useSelector<RootState, boolean>(
    state => state.auth.loggedIn
  );
  const dispatch = useDispatch();

  return (
    <AppBar position="fixed" color="secondary">
      <Toolbar>
        <Link to="/" className={classes.icon}>
          <img src={logo} alt="Kiuru Logo" />
        </Link>
        {!loggedIn && (
          <div className={classes.menu}>
            <Button color="primary" component={Link} to="/signup">
              <Typography variant="h6">Sign up</Typography>
            </Button>
            <Button color="primary" component={Link} to="/login">
              <Typography variant="h6">Log in</Typography>
            </Button>
          </div>
        )}
        {loggedIn && (
          <div className={classes.menu}>
            <Button
              color="primary"
              component={Link}
              to="/"
              onClick={() => dispatch(logout())}
            >
              <Typography variant="h6">Log out</Typography>
            </Button>
          </div>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default TopNav;
