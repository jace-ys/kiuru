import React from "react";

import { Hidden } from "@material-ui/core";

const Display: React.FC = props => {
  return <Hidden xsDown>{props.children}</Hidden>;
};

export default Display;
