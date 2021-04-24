import React from "react";
import Layout from "@theme/Layout";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Hero from "./components/Hero/hero";
import Deliver from "./components/Deliver/Deliver";
import WhySysl from "./components/WhySysl/WhySysl";

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="The powerful system specification language"
    >
      <Hero />
      <main>
        <Deliver />
        <WhySysl />
      </main>
    </Layout>
  );
}

export default Home;
