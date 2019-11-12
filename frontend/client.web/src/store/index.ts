import { combineReducers, configureStore } from "redux-starter-kit";

import authReducer from "./ducks/auth";

const rootReducer = combineReducers({
  auth: authReducer
});

export type RootState = ReturnType<typeof rootReducer>;

const store = configureStore({
  reducer: rootReducer
});

export default store;
