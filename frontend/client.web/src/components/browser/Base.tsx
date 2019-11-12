import React from "react";

import { Hidden } from "@material-ui/core";

import TopNav from "./TopNav";

const Base: React.FC = () => {
  return (
    <Hidden xsDown>
      <TopNav />
    </Hidden>
  );
};

export default Base;
