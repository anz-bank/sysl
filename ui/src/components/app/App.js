import React from "react";
import { HashRouter } from "react-router-dom";
import AppMain from "./AppMain";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      screenWidth: 0,
      screenHeight: 0
    };
    this.updateWindowDimensions = this.updateWindowDimensions.bind(this);
  }

  componentDidMount() {
    this.updateWindowDimensions();
    window.addEventListener("resize", this.updateWindowDimensions);
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.updateWindowDimensions);
  }

  updateWindowDimensions() {
    const { state } = this;
    state.screenWidth = window.innerWidth;
    state.screenHeight = window.innerHeight;
    this.setState(state);
  }

  render() {
    return (
      <HashRouter>
        <AppMain />
      </HashRouter>
    );
  }
}

export default App;
