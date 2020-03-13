import React from "react";
import { styled } from "@material-ui/core/styles";
import { Route } from "react-router-dom";
import { Container } from "@material-ui/core";
import Box from "@material-ui/core/Box";
import Routes from "../router/Routes";

const StyledContainer = styled(Container)({
  backgroundColor: "transparent",
  paddingTop: 65,
  paddingLeft: 40,
  paddingRight: 40
});

const StyledBox = styled(Box)({
  display: "flex",
  padding: 0
});

const AppMain = () => {
  return (
    <>
      <StyledBox>
        <StyledContainer>
          <Route render={({ location }) => <Routes location={location} />} />
        </StyledContainer>
      </StyledBox>
    </>
  );
};

export default AppMain;
