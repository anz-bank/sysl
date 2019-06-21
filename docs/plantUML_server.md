# Installing PlantUML on Apache Tomcat

## Mac OSX guide

You want to install apache tomcat with homebrew

    brew install tomcat

Now go to http://plantuml.com/download and look for the section `Java J2EE WAR File` to download [`plantuml.war`](http://sourceforge.net/projects/plantuml/files/plantuml.war/download)

Once you have installed apache tomcat copy your 'plantuml.war' into your tomcat server:

    cp plantuml.war /usr/local/Cellar/tomcat/<catalina version>/libexec/webapps

You can find your catalina version under `CATALINA_BASE`:

    catalina version

Now start your tomcat server

    catalina start

In a web browser of `your choice, you can navigate to `localhost:8080/plantuml` for your plantuml server!
