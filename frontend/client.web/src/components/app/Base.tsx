import React from "react";
import { useSelector } from "react-redux";

import { Hidden } from "@material-ui/core";

import Auth from "./Auth";
import BottomNav from "./BottomNav";

import { RootState } from "../../store";

const Base: React.FC = () => {
  const loggedIn = useSelector<RootState, boolean>(
    state => state.auth.loggedIn
  );

  return (
    <Hidden smUp>
      {!loggedIn && <Auth />}
      {loggedIn && <BottomNav />}
    </Hidden>
  );
};

export default Base;
