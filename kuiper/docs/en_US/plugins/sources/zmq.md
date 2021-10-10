# Zmq Source

The source will subscribe to a Zero Mq topic to import the messages into kuiper

## Compile & deploy plugin

```shell
# cd $kuiper_src
# go build --buildmode=plugin -o plugins/sources/Zmq.so plugins/sources/zmq.go
# cp plugins/sources/Zmq.so $kuiper_install/plugins/sources
```

Restart the Kuiper server to activate the plugin.

## Configuration

The configuration for this source is ``$kuiper/etc/sources/zmq.yaml``. The format is as below:

```yaml
#Global Zmq configurations
default:
  server: tcp://192.168.2.2:5563  
test:
  server: tcp://127.0.0.1:5563
```
### Global configurations

Use can specify the global zmq source settings here. The configuration items specified in ``default`` section will be taken as default settings for the source when connects to Zero Mq.

### server

The url of the Zero Mq server that the source will subscribe to.

## Override the default settings

If you have a specific connection that need to overwrite the default settings, you can create a customized section. In the previous sample, we create a specific setting named with ``test``.  Then you can specify the configuration with option ``CONF_KEY`` when creating the stream definition (see [stream specs](../../sqls/streams.md) for more info).

## Sample usage

```
demo (
		...
	) WITH (DATASOURCE="demo", FORMAT="JSON", CONF_KEY="test", TYPE="zmq");
```

The configuration keys "test" will be used. The Zero Mq topic to subscribe is "demo" as specified in the ``DATASOURCE``.

