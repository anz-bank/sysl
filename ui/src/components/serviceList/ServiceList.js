import React, { useState } from "react";
import ServiceListView from "./ServiceListView";

function getData() {
  return fetch("./data/services.json").then(res => res.json());
}

function ServiceList() {
  const [data, setData] = useState([]);
  // Todo: Fix this so that it's only called on load
  getData().then(newData => setData(newData));
  return <ServiceListView items={data} />;
}

export default ServiceList;
