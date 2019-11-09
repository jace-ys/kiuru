import React from "react";
import { Link } from "react-router-dom";

import { BottomNavigation, BottomNavigationAction } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import {
  AccountCircleRounded,
  NotificationsRounded,
  PeopleAltRounded,
  StarRounded
} from "@material-ui/icons";

const useStyles = makeStyles({
  root: {
    position: "fixed",
    width: "100%",
    bottom: 0,
    borderTop: "0.2px solid #e0e0e0",
    boxShadow: "0px -0.5px 5px -1px #e0e0e0",
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
    <BottomNavigation
      value={value}
      onChange={handleChange}
      className={classes.root}
    >
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
  );
};

export default BottomNav;
