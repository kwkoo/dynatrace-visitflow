# Dynatrace Visit Flow Diagram

## Overview

* This application consumes the user session export feed from Dynatrace and renders a sankey diagram (using <https://github.com/d3/d3-sankey>) based on the user session data.
* It maintains a graph of user actions in memory. It does not persist the graph anywhere. This means that the graph is lost if the application is stopped.
* The easiest way to run this is by using Docker - build the image with the included `buildimage` script, and execute the image with the `run` script.
* By default, the server starts up a listener on port 8080. It accepts user actions from all applications. If you wish to capture user actions for a specific application, you can modify the run script to include the `-application` argument when starting the container.
* If you are running this from a non Internet-accessible machine, you can use [ngrok](https://ngrok.com/) to create a tunnel between your machine and the Internet. You can then configure Dynatrace to export user sessions to the `ngrok.io` host.
* Visits are only exported from Dynatrace once they are completed. This means that you may only see the session after a 30-minute timeout.

## Useful URLs

Here are a list of useful URLs (assuming you're running this on your local machine):

* Web interface - <http://localhost:8080>
* Dump graph to browser - <http://localhost:8080/dump>
* Generate JSON - <http://localhost:8080/api>

