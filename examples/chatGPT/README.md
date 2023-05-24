## Sample Plugin for ChatGPT-4

This sample intends to be a plugin to generate and display maps in svg directly in ChatGPT

### Prerequisites

- A GPT Plus subscription
- The Go toolchain

### Configuration

There is two configurations possible:

- [the plugin manifest](https://platform.openai.com/docs/plugins/getting-started/plugin-manifest). This generates the `ai-plugin.json` served via `/.well-known/ai-plugin.json`
- the configuration of the plugin, listen ports and so on that is made through env variables.

#### Configuration of the manifest

To change the manifest, you should edit the file `wellknown.cue`. This file is expressed in [CUE](https://cuelang.org).
A validation of the content is performed when the plugin starts.

The file is transformed into JSON at runtime

#### Configuration of the plugin

This application is configured via the environment. The following environment
variables can be used:

```text 
KEY                        TYPE      DEFAULT           REQUIRED    DESCRIPTION
WTG_CHATGPT_LISTEN_ADDR    String    localhost:3333    true        Host to connect to, or ngrok to use tunneling
```

So by default the services listens on `localhost:3333`

If you specify `ngrok` for the variable ``WTG_CHATGPT_LISTEN_ADDR` and you have the ngrok token, a tunnel should be setup.
The listen address is automatically changed in the various files such as openai.yaml or ai-plugin.json.

### Configure chatGPT

Once you have the server up and running, you can go into ChatGPT and try "_Develop your own plugin_" and enter `localhost:3333` in the URL.
