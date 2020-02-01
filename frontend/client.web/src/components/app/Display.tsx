import React from "react";

import { Hidden } from "@material-ui/core";

const Display: React.FC = props => {
  return <Hidden smUp>{props.children}</Hidden>;
};

export default Display;
