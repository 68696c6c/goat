# GOAT
Goat is a toolkit for creating JSON-speaking RESTful APIs in Go. 

Goat aims to streamline development by offering a tightly focused feature set, convention-over-configuration philosophies, and out-of-the-box devops tooling.

Goat is *not* made for building traditional webapps with a front-end.  Instead, it focuses on helping developers create headless webservices that can be used as the backend for server-side rendered static sites or mobile-apps.

# Roadmap

Planned features include:

- A global CLI tool for for scaffolding projects (including out-of-the-box Docker support and Makefiles)
- Full migration suite
- Generator scripts for Cobra commands, repositories, models, handlers, migrations, etc.
- Out-of-the-box AWS Cloudformation support via generated templates and bootstrapping scripts, including Parameter Store management of secrets and ECS clusters.
- Oauth support
- Out-of-the-box Cobra utility commands 
- Support for more additional database engines (currently only supports MySQL-compatible databases)
- Front end integrations via Horns, a comprehensive, fully-themable, open-source UI framework made for Gatsby and React Emotion with React Native support   


# Demo
*This section is out of date*

A basic demo app can be found in the `demo` directory.  To run it, change 
into the directory and run:

```
dep ensure
./run.sh
```

Running the demo app will print some information on how to use Goat.
