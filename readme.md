## S3 JSON editor and visualizer

This is a tool that lets you modify only JSON files on S3.

You can modify either RFC compliant JSON or Athena compliant JSON (one line per JSON without any comma separator).

## How to build
- Launch the update_front_end.sh script to create the needed front end files.
- Go into server/build/index.html and modify as follows:

FROM:
```
```
<script type="module" crossorigin src="/assets/index.cbe63ed5.js"></script>
    <link rel="stylesheet" href="/assets/index.8050d0fe.css">
```
TO:

```
<script type="module" crossorigin src="./assets/index.cbe63ed5.js"></script>
    <link rel="stylesheet" href="./assets/index.8050d0fe.css">
```

- Launch the deploy.sh script

## How to run
Just launch the executable

If on windows webview requires some .DLL that you need to put in the same directory of the executable
Please Refer to: https://github.com/webview/webview/issues/404

## Configuration
As of now only the server port can be configured. Put the ```config.json``` file near the executable to change it. If the file is not present the webserver will start on the 3005 port




