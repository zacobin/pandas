# Stream specs 

## Data types

In Kuiper, each column or an expression has a related data type. A data type describes (and constrains) the set of values that a column of that type can hold or an expression of that type can produce.

Below is the list of data types supported.

| #    | Data type | Description                                                  |
| ---- | --------- | ------------------------------------------------------------ |
| 1    | bigint    |                                                              |
| 2    | float     |                                                              |
| 3    | string    |                                                              |
| 4    | datetime  | Not support.                                                 |
| 5    | boolean   |                                                              |
| 6    | array     | The array type, can be any simple types or struct type (#1 - #5, and #7). |
| 7    | struct    | The complex type.                                            |

## Language definitions

```sql
CREATE STREAM   
    stream_name   
    ( column_name <data_type> [ ,...n ] )
    WITH ( property_name = expression [, ...] );
```

**The supported property names.**

| Property name | Optional | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| DATASOURCE | false    | The value is determined by source type. The topic names list if it's a MQTT data source. Please refer to related document for other sources. |
| FORMAT        | true | The data format, currently the value can only be "JSON". |
| KEY           | true     | Reserved key, currently the field is not used. It will be used for GROUP BY statements. |
| TYPE     | true | The source type, if not specified, the value is "mqtt". |
| StrictValidation     | true | To control validation behavior of message field against stream schema. See [StrictValidation](#StrictValidation) for more info. |
| CONF_KEY | true | If additional configuration items are requied to be configured, then specify the config key here. See [MQTT stream](../rules/sources/mqtt.md) for more info. |

**Example 1,**

```sql
my_stream 
  (id bigint, name string, score float)
WITH ( datasource = "topic/temperature", FORMAT = "json", KEY = "id");
```

The stream will subscribe to MQTT topic ``topic/temperature``, the server connection uses ``servers`` key of ``default`` section in configuration file ``$kuiper/etc/mqtt_source.yaml``. 

- See [MQTT source](../rules/sources/mqtt.md) for more info.

**Example 2,**

```sql
demo (
		USERID BIGINT,
		FIRST_NAME STRING,
		LAST_NAME STRING,
		NICKNAMES ARRAY(STRING),
		Gender BOOLEAN,
		ADDRESS STRUCT(STREET_NAME STRING, NUMBER BIGINT),
	) WITH (DATASOURCE="test/", FORMAT="JSON", KEY="USERID", CONF_KEY="demo");
```

The stream will subscribe to MQTT topic ``test/``, the server connection uses settings of ``demo`` section in configuration file ``$kuiper/etc/mqtt_source.yaml``. 

- See [MQTT source](../rules/sources/mqtt.md) for more info.

- See [rules and streams CLI docs](../cli/overview.md) for more information of rules & streams management.

### StrictValidation

```
The value of StrictValidation can be true or false.
1) True: Drop the message if the message  is not satisfy with the stream definition.
2) False: Keep the message, but fill the missing field with default empty value.

bigint: 0
float: 0.0
string: ""
datetime: (NOT support yet)
boolean: false
array: zero length array
struct: null value
```

### Schema-less stream
If the data type of the stream is unknown or varying, we can define it without the fields. This is called schema-less. It is defined by leaving the fields empty.
```sql
schemaless_stream 
  ()
WITH ( datasource = "topic/temperature", FORMAT = "json", KEY = "id");
```

Schema-less stream field data type will be determined at runtime. If the field is used in an incompatible clause, a runtime error will be thrown and send to the sink. For example, ``where temperature > 30``. Once a temperature is not a number, an error will be sent to the sink.

See [Query languange element](query_language_elements.md) for more inforamtion of SQL language.

