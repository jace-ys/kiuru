import { createSlice } from "redux-starter-kit";

export interface AuthState {
  loggedIn: boolean;
}

let initialState: AuthState = {
  loggedIn: false
};

const authSlice = createSlice({
  name: "auth",
  initialState: initialState,
  reducers: {
    login: state => {
      state.loggedIn = true;
    },
    logout: state => {
      state.loggedIn = false;
    }
  }
});

export const { login, logout } = authSlice.actions;

export default authSlice.reducer;
