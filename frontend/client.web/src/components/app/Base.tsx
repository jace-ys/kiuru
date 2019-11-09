import React, { useState } from "react";

import Auth from "./Auth";
import BottomNav from "./BottomNav";
import TopNav from "./TopNav";

const Base: React.FC = () => {
  const [loggedIn, setLoggedIn] = useState(false);

  return (
    <div>
      {!loggedIn && <Auth loggedIn={loggedIn} setLoggedIn={setLoggedIn} />}
      {loggedIn && (
        <div>
          <TopNav title="Discover" />
          <BottomNav />
        </div>
      )}
    </div>
  );
};

export default Base;
