import React from "react";
import { Link } from "react-router-dom";
import { useDispatch } from "react-redux";

import { Avatar, Button, Grid, TextField, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import { LockOutlined } from "@material-ui/icons";

import Base from "../components/Base";

import { login } from "../store/ducks/auth";

const useStyles = makeStyles(theme => ({
  image: {
    backgroundImage:
      "url(https://images.squarespace-cdn.com/content/v1/587271d7ff7c50c708f3e44b/1578787510843-P4GBYL9307QO2HR2WCZ5/ke17ZwdGBToddI8pDm48kOeC2_vIwtkNO4KLfB1JIsJ7gQa3H78H3Y0txjaiv_0fDoOvxcdMmMKkDsyUqMSsMWxHk725yiiHCCLfrh8O1z5QPOohDIaIeljMHgDF5CVlOqpeNLcJ80NK65_fV7S1UedT6MCuDrG0_6iPwXLGF1669zXNhvZ0Gt3DqtjtSHkNlcTmcwU-Ed_fLjLxuPb0KQ/HaifossWaterfallIceland.jpg)",
    backgroundPosition: "center",
    backgroundRepeat: "no-repeat",
    backgroundSize: "cover"
  },
  formPanel: {
    padding: theme.spacing(4),
    display: "flex",
    flexDirection: "column",
    alignItems: "center"
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.primary.main
  },
  title: {
    margin: theme.spacing(0.5)
  },
  submit: {
    margin: theme.spacing(3, 0, 2)
  }
}));

interface Field {
  id: string;
  label: string;
  name: string;
  autoComplete: string;
}

const Login: React.FC = () => {
  const classes = useStyles();
  const dispatch = useDispatch();

  const fields: Field[] = [
    {
      id: "email",
      label: "Username / Email",
      name: "email",
      autoComplete: "email"
    },
    {
      id: "password",
      label: "Password",
      name: "password",
      autoComplete: "password"
    }
  ];

  return (
    <Base browser>
      <Grid container>
        <Grid item xs={false} sm={6} md={7} className={classes.image} />
        <Grid item xs={12} sm={6} md={5}>
          <div className={classes.formPanel}>
            <Avatar className={classes.avatar}>
              <LockOutlined />
            </Avatar>
            <Typography component="h1" variant="h4" className={classes.title}>
              Log in
            </Typography>
            <Typography component="h1" variant="subtitle1">
              Don't have an account? <Link to="/signup">Sign up</Link>
            </Typography>
            <form noValidate>
              <Grid container>
                {fields.map((field, index) => {
                  return (
                    <TextField
                      key={index}
                      variant="outlined"
                      margin="normal"
                      required
                      fullWidth
                      id={field.id}
                      label={field.label}
                      name={field.name}
                      autoComplete={field.autoComplete}
                      autoFocus={index ? false : true}
                    />
                  );
                })}
                <Button
                  onClick={() => dispatch(login())}
                  component={Link}
                  to="/discover"
                  variant="contained"
                  color="primary"
                  size="large"
                  fullWidth
                  className={classes.submit}
                >
                  Log in
                </Button>
                <Grid container justify="center">
                  <Typography component="h1" variant="subtitle2">
                    <Link to="#">Forgot your password?</Link>
                  </Typography>
                </Grid>
              </Grid>
            </form>
          </div>
        </Grid>
      </Grid>
    </Base>
  );
};

export default Login;
