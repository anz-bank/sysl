import React from "react";
import classnames from "classnames";
import useBaseUrl from "@docusaurus/useBaseUrl";
import styles from "./styles.module.css";

const features = [
  {
    title: <>Diagrams</>,
    imageUrl: "img/icon/diagram.svg",
    description: (
      <>
        Sysl generates data, integration and sequence diagrams from your specs.
      </>
    ),
  },
  {
    title: <>Code</>,
    imageUrl: "img/icon/code.svg",
    description: (
      <>
        Sysl generates Go applications out of the box: client and server, REST
        and gRPC.
      </>
    ),
  },
  {
    title: <>Extensions</>,
    imageUrl: "img/icon/plug.svg",
    description: (
      <>
        Sysl supports an ecosystem of tools - generate anything else that you
        need.
      </>
    ),
  },
  {
    title: <>Synchronised</>,
    imageUrl: "img/icon/sync.svg",
    description: (
      <>
        Your Sysl specifications become your single source of truth - no more
        stale designs.
      </>
    ),
  },
  {
    title: <>Cross-platform</>,
    imageUrl: "img/icon/computer.svg",
    description: <>Sysl works on Windows, MacOS and Linux.</>,
  },
  {
    title: <>Open Source</>,
    imageUrl: "img/icon/lock_open.svg",
    description: (
      <>
        Sysl is licensed under Apache 2.0, free for personal or commercial use.
      </>
    ),
  },
];

function Feature({ imageUrl, title, description }) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={classnames("col col--4", styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <object
            type="image/svg+xml"
            className={styles.featureImage}
            data={imgUrl}
          ></object>
        </div>
      )}
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  );
}

function Features() {
  return (
    <React.Fragment>
      {features && features.length > 0 && (
        <section className={styles.features}>
          <div className="container">
            <div className="row">
              {features.map((props, idx) => (
                <Feature key={idx} {...props} />
              ))}
            </div>
          </div>
        </section>
      )}
    </React.Fragment>
  );
}

export default Features;
