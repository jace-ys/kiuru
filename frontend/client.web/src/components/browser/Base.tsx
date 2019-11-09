import React from "react";

import TopNav from "./TopNav";

const Base: React.FC = () => {
  const loggedIn = false;

  return (
    <div>
      <TopNav loggedIn={loggedIn} />
    </div>
  );
};

export default Base;
