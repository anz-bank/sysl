import { Route, Switch } from "react-router";
import React from "react";
import posed, { PoseGroup } from "react-pose";
import ServiceList from "components/serviceList/ServiceList";

const RoutesContainer = posed.div({
  enter: { opacity: 1, beforeChildren: true },
  exit: { opacity: 0 }
});

export default function Routes(props) {
  return (
    <PoseGroup>
      <h1 key="title">Sysl UI</h1>
      <RoutesContainer key={props.location.pathname}>
        <Switch location={props.location}>
          <Route path="/" component={ServiceList} />
        </Switch>
      </RoutesContainer>
    </PoseGroup>
  );
}
