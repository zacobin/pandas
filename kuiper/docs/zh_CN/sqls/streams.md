# 流规格

## 数据类型



在 Kuiper 中，每个列或表达式都有一个相关的数据类型。 数据类型描述（约束）该类型的列可以容纳的一组值或该类型可以产生的表达式。

以下是支持的数据类型的列表。

| #    | 数据类型 | 说明                                                   |
| ---- | -------- | ------------------------------------------------------ |
| 1    | bigint   |                                                        |
| 2    | float    |                                                        |
| 3    | string   |                                                        |
| 4    | datetime | 不支持                                                 |
| 5    | boolean  |                                                        |
| 6    | array    | 数组类型可以是任何简单类型或结构类型（＃1-＃5和＃7）。 |
| 7    | struct   | 复杂类型                                               |

## 语言定义

```sql
CREATE STREAM   
    stream_name   
    ( column_name <data_type> [ ,...n ] )
    WITH ( property_name = expression [, ...] );
```

**支持的属性名称**

| 属性名称 | 可选 | 说明                                              |
| ------------- | -------- | ------------------------------------------------------------ |
| DATASOURCE | 否   | 取决于不同的源类型；如果是 MQTT 源，则为 MQTT 数据源主题名；其它源请参考相关的文档。 |
| FORMAT        | 是      | 传入的数据类型，目前只支持 "JSON"。 |
| KEY           | 是    | 保留配置，当前未使用该字段。 它将用于 GROUP BY 语句。 |
| TYPE    | 是      | 源类型，如未指定，值为 "mqtt"。 |
| StrictValidation     | 是  | 针对流模式控制消息字段的验证行为。 有关更多信息，请参见 [StrictValidation](#StrictValidation) |
| CONF_KEY | 是 | 如果需要配置其他配置项，请在此处指定 config 键。 有关更多信息，请参见 [MQTT stream](../rules/sources/mqtt.md) 。 |

**示例1**

```sql
my_stream 
  (id bigint, name string, score float)
WITH ( datasource = "topic/temperature", FORMAT = "json", KEY = "id");
```

该流将订阅 MQTT 主题`topic/temperature`，服务器连接使用配置文件`$kuiper/etc/mqtt_source.yaml` 中默认部分的 servers 键。

- 有关更多信息，请参见 [MQTT source](../rules/sources/mqtt.md) 

**示例2**

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

 流将订阅 MQTT 主题 `test/`，服务器连接使用配置文件`$kuiper/etc/mqtt_source.yaml` 中 demo 部分的设置。

- 有关更多信息，请参见 [MQTT source](../rules/sources/mqtt.md) 

- 有关规则和流管理的更多信息，请参见 [规则和流 CLI docs](../cli/overview.md) 

### StrictValidation

```
StrictValidation 的值可以为 true 或 false。
1）True：如果消息不符合流定义，则删除消息。
2）False：保留消息，但用默认的空值填充缺少的字段。

bigint: 0
float: 0.0
string: ""
datetime: (NOT support yet)
boolean: false
array: zero length array
struct: null value
```

### Schema-less 流

如果流的数据类型未知或不同，我们可以不使用字段来定义它。 这称为 schema-less。 通过将字段设置为空来定义它。

```sql
schemaless_stream 
  ()
WITH ( datasource = "topic/temperature", FORMAT = "json", KEY = "id");
```

Schema-less 流字段数据类型将在运行时确定。 如果在不兼容子句中使用该字段，则会抛出运行时错误并将其发送到目标。 例如，`where temperature > 30`。 一旦温度不是数字，将错误发送到目标。

有关 SQL 语言的更多信息，请参见 [查询语言元素](query_language_elements.md) 。

