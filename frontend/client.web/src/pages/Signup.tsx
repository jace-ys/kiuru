import React from "react";
import { Link } from "react-router-dom";
import { useDispatch } from "react-redux";

import { Avatar, Button, Grid, TextField, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import { PersonAddOutlined } from "@material-ui/icons";

import Base from "../components/Base";

import { login } from "../store/ducks/auth";

const useStyles = makeStyles(theme => ({
  image: {
    backgroundImage:
      "url(https://images.squarespace-cdn.com/content/v1/587271d7ff7c50c708f3e44b/1568586724159-WH56CI2XN0M0FNCE93EC/ke17ZwdGBToddI8pDm48kJKo3YTR7zgUvInmXMbZ6zZ7gQa3H78H3Y0txjaiv_0fDoOvxcdMmMKkDsyUqMSsMWxHk725yiiHCCLfrh8O1z4YTzHvnKhyp6Da-NYroOW3ZGjoBKy3azqku80C789l0geeCvn1f36QDdcifB7yxGjTk-SMFplgtEhJ5kBshkhu5q5viBDDnY2i_eu2ZnquSA/_DSC6453.jpg)",
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
              <PersonAddOutlined />
            </Avatar>
            <Typography component="h1" variant="h4" className={classes.title}>
              Sign up
            </Typography>
            <Typography component="h1" variant="body1">
              Already have an account? <Link to="/login">Log in</Link>
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
