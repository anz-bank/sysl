import React from "react";
import classnames from "classnames";
import styles from "./styles.module.css";
import deliveryStyles from "../Deliver/styles.module.css";

const badges = [
  {
    title: <>Guaranteed Consistency</>,
    imageUrl: "img/icon/feature-consistency.svg",
    description: (
      <>
        Sysl is built for change. With Sysl, your specifications become the
        single source of truth. All the related system specifications, code,
        diagrams, documentations keep the same version - no more stale designs.
        Also, in unified good style.
      </>
    ),
  },
  {
    title: <>Quality Built-in</>,
    imageUrl: "img/icon/feature-quality.svg",
    description: (
      <>
        With Sysl, you can generate high-quality code, diagrams and
        documentation which align to your style guide and best practices. Embed
        security in from the start.
      </>
    ),
  },
  {
    title: <>Flexible</>,
    imageUrl: "img/icon/feature-flexible.svg",
    description: (
      <>
        Sysl's open-ended architecture makes it possible to generate code for
        any programming language, with out-of-the-box support for Go.
      </>
    ),
  },
];

function WhySysl() {
  return (
    <React.Fragment>
      <div id="why" className={styles.why__sysl__title}>
        <a href="#why">
          <h1
            className={classnames(
              "text--primary text--center",
              deliveryStyles.deliver__title
            )}
          >
            Why Sysl
          </h1>
        </a>
      </div>
      <div className={styles.why__sysl__curve}></div>
      <div className={styles.why__sysl}>
        <div className={styles.why__sysl__badge__group}>
          {badges &&
            badges.length > 0 &&
            badges.map((badge, idx) => {
              return (
                <div className={styles.why__sysl__badge} key={idx}>
                  <img src={badge.imageUrl}></img>
                  <h1 className="text--center">{badge.title}</h1>
                  <span>{badge.description}</span>
                </div>
              );
            })}
        </div>
      </div>
    </React.Fragment>
  );
}

export default WhySysl;
