import React from "react";
import { Link } from "react-router-dom";
import { useDispatch } from "react-redux";

import { Avatar, Button, Grid, TextField, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import { PersonOutlined } from "@material-ui/icons";

import Base from "../components/Base";

import { login } from "../store/ducks/auth";

import signup from "../assets/images/signup.jpg";

const useStyles = makeStyles(theme => ({
  image: {
    backgroundImage: `url(${signup})`,
    backgroundPosition: "center",
    backgroundRepeat: "no-repeat",
    backgroundSize: "cover"
  },
  formPanel: {
    padding: theme.spacing(4),
    textAlign: "center"
  },
  avatar: {
    margin: theme.spacing(1),
    marginLeft: "auto",
    marginRight: "auto",
    backgroundColor: theme.palette.primary.main
  },
  text: {
    margin: theme.spacing(0.5),
    "& a": {
      color: theme.palette.action.hover
    }
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
      id: "name",
      label: "Name",
      name: "name",
      autoComplete: "name"
    },
    {
      id: "username",
      label: "Username",
      name: "username",
      autoComplete: "username"
    },
    {
      id: "email",
      label: "Email",
      name: "email",
      autoComplete: "email"
    },
    {
      id: "password",
      label: "Password",
      name: "password",
      autoComplete: "current-password"
    }
  ];

  return (
    <Base browser>
      <Grid container>
        <Grid item xs={false} sm={6} md={7} className={classes.image} />
        <Grid item xs={12} sm={6} md={5}>
          <div className={classes.formPanel}>
            <Avatar className={classes.avatar}>
              <PersonOutlined />
            </Avatar>
            <Typography variant="h4" className={classes.text}>
              Sign up
              <Typography variant="body1">
                Already have an account? <Link to="/login">Log in</Link>
              </Typography>
            </Typography>
            <form noValidate>
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
                Sign up
              </Button>
            </form>
          </div>
        </Grid>
      </Grid>
    </Base>
  );
};

export default Login;
