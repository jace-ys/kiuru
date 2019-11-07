import React from "react";

import BottomNav from "./BottomNav";
import TopNav from "./TopNav";

const AppBase: React.FC = () => {
  return (
    <div className="AppBase">
      <TopNav title="Discover" />
      <BottomNav />
    </div>
  );
};

export default AppBase;
