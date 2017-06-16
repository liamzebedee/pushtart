
Pushtart - The worlds easiest PaaS. [![Build Status](https://travis-ci.org/twitchyliquid64/pushtart.svg?branch=master)](https://travis-ci.org/twitchyliquid64/pushtart)
=======================================

Need to deploy a project? Sick of learning new deployment strategies, configurations, custom pipelines/protocols?

Pushtart is the Unix of deployment:
 - Deployment is pushing to a remote Git repo
 - Setup/running/building is a Bash script -- no custom DSL/language
 - Administering deployments (starting, stopping, environment vars, etc.) is transported over SSH
 
Let's see how simple it is...

## An example
 1. You have your Node/Python/Brainfuck project on your local computer - /home/liam/project.
 2. We need to write the startup script for deployment. It's going to be called `startup.sh`:
```bash
PORT=80 node server.js
```
 3. Install Pushtart
 4. Run `pushtart`. This will ask you for the name of your project (by default it's the parent directory), your SSH login for your server, and the name of your Git branch for deployment.
 5. Pushtart will automagically create this Git branch, install itself on your server, and then deploy.

To redeploy, it's as simple as:
 1. Commit your changes
 2. `git push pushtart_server master`

Pushtart will also save your configuration to tartconfig, so you can simply set this up on any other server.

## What else?
If it was this simple, we'd be using Git post-deploy hooks. Pushtart offers a bit more...

 - Environment variables
 - Automatic restart
 - STDOUT/STDERR log management
 - DNS server
 - HTTP Proxy
