import React from "react";
import classnames from "classnames";
import styles from "./styles.module.css";
import Link from "@docusaurus/Link";
import useBaseUrl from "@docusaurus/useBaseUrl";

function Deliver() {
  return (
    <React.Fragment>
      <div id="deliver" className={styles.deliver__title__section}>
        <a href="#deliver">
          <h1
            className={classnames(
              "text--primary text--center",
              styles.deliver__title
            )}
          >
            What Sysl Delivers
          </h1>
        </a>
        <p className={styles.deliver__desc}>
          Generate design diagrams in seconds. Create new services in minutes.
          <br />
          Take those specifications and get all of Sysl's other benefits for
          free.
        </p>
      </div>
      <div className={styles.deliver__section}>
        <div className={styles.deliver__body}>
          <div className={styles.deliver__left}>
            <p className={styles.deliver__left__title}>Sources</p>

            <div className={styles.deliver__button__group__left}>
              <Link
                to={useBaseUrl("docs/lang-spec")}
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Sysl
              </Link>
              <Link
                to={useBaseUrl("docs/cmd/cmd-import")}
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                OpenAPI2/OpenAPI3
              </Link>
              <Link
                to={useBaseUrl("docs/cmd/cmd-import")}
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Protobuf
              </Link>
              <Link
                to={useBaseUrl("docs/cmd/cmd-import")}
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                XSD
              </Link>
            </div>
            <div className={styles.deliver__arrow__left}>
              <img src="img/icon/arrow_in_group.svg"></img>
            </div>
          </div>
          <div className={styles.deliver__middle}>
            <img src="img/icon/sysl_circle.svg"></img>
            <h1
              className={classnames(
                "text--primary text--center",
                styles.deliver__middle__title
              )}
            >
              Sysl
            </h1>
          </div>

          <div className={styles.deliver__right}>
            <p className={styles.deliver__right__title}>Outcomes</p>
            <div className={styles.deliver__arrow__right__group}>
              <img src="img/icon/arrow_right.svg"></img>
              <img src="img/icon/arrow_right.svg"></img>
              <img src="img/icon/arrow_right.svg"></img>
              <img src="img/icon/arrow_right.svg"></img>
            </div>
            <div className={styles.deliver__button__group__right}>
              <Link
                to={useBaseUrl("docs/gen-diagram")}
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Diagram
              </Link>
              <Link
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Code
              </Link>
              <Link
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Documentation
              </Link>
              <Link
                className={classnames(
                  "button button--primary button--lg",
                  styles.button__round
                )}
              >
                Tests
              </Link>
            </div>
          </div>
        </div>
      </div>
    </React.Fragment>
  );
}

export default Deliver;
