import React from "react";
import { Link } from "react-router-dom";

import {
  BottomNavigation,
  BottomNavigationAction,
  Paper
} from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import {
  AccountCircleRounded,
  NotificationsRounded,
  PeopleAltRounded,
  StarRounded
} from "@material-ui/icons";

const useStyles = makeStyles({
  root: {
    bottom: 0,
    position: "fixed",
    width: "100vw",
    "& .MuiBottomNavigationAction-root": {
      padding: 0
    }
  }
});

const BottomNav: React.FC = () => {
  const classes = useStyles();
  const [value, setValue] = React.useState("discover");
  const handleChange = (event: React.ChangeEvent<{}>, newValue: string) => {
    setValue(newValue);
  };

  return (
    <Paper elevation={6} className={classes.root}>
      <BottomNavigation value={value} onChange={handleChange}>
        <BottomNavigationAction
          component={Link}
          to="/discover"
          value="discover"
          icon={<StarRounded />}
        />
        <BottomNavigationAction
          component={Link}
          to="/connect"
          value="connect"
          icon={<PeopleAltRounded />}
        />
        <BottomNavigationAction
          component={Link}
          to="/notifications"
          value="notifications"
          icon={<NotificationsRounded />}
        />
        <BottomNavigationAction
          component={Link}
          to="/profile"
          value="profile"
          icon={<AccountCircleRounded />}
        />
      </BottomNavigation>
    </Paper>
  );
};

export default BottomNav;
