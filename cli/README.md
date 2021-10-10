# Pandas CLI
## Build
From the project root:
```
make cli
```

## Usage
### Service
#### Get the version of Pandas services
```
pandas-cli version
```

### Users management
#### Create User
```
pandas-cli users create <user_email> <user_password>
```

#### Login User
```
pandas-cli users token <user_email> <user_password>
```

#### Retrieve User
```
pandas-cli users get <user_auth_token>
```

#### Update User Metadata
```
pandas-cli users update '{"key1":"value1", "key2":"value2"}' <user_auth_token>
```

#### Update User Password
```
pandas-cli users password <old_password> <password> <user_auth_token>
```

### System Provisioning
#### Create Thing (type Device)
```
pandas-cli things create '{"name":"myDevice"}' <user_auth_token>
```

#### Create Thing (type Application)
```
pandas-cli things create '{"name":"myDevice"}' <user_auth_token>
```

#### Update Thing
```
pandas-cli things update '{"id":"<thing_id>", "name":"myNewName"}' <user_auth_token>
```

#### Remove Thing
```
pandas-cli things delete <thing_id> <user_auth_token>
```

#### Retrieve a subset list of provisioned Things
```
pandas-cli things get all --offset=1 --limit=5 <user_auth_token>
```

#### Retrieve Thing By ID
```
pandas-cli things get <thing_id> <user_auth_token>
```

#### Create Channel
```
pandas-cli channels create '{"name":"myChannel"}' <user_auth_token>
```

#### Update Channel
```
pandas-cli channels update '{"id":"<channel_id>","name":"myNewName"}' <user_auth_token>

```
#### Remove Channel
```
pandas-cli channels delete <channel_id> <user_auth_token>
```

#### Retrieve a subset list of provisioned Channels
```
pandas-cli channels get all --offset=1 --limit=5 <user_auth_token>
```

#### Retrieve Channel By ID
```
pandas-cli channels get <channel_id> <user_auth_token>
```

### Access control
#### Connect Thing to Channel
```
pandas-cli things connect <thing_id> <channel_id> <user_auth_token>
```

#### Disconnect Thing from Channel
```
pandas-cli things disconnect <thing_id> <channel_id> <user_auth_token>

```

#### Retrieve a subset list of Channels connected to Thing
```
pandas-cli things connections <thing_id> <user_auth_token>
```

#### Retrieve a subset list of Things connected to Channel
```
pandas-cli channels connections <channel_id> <user_auth_token>
```

### Messaging
#### Send a message over HTTP
```
pandas-cli messages send <channel_id> '[{"bn":"Dev1","n":"temp","v":20}, {"n":"hum","v":40}, {"bn":"Dev2", "n":"temp","v":20}, {"n":"hum","v":40}]' <thing_auth_token>
```

#### Read messages over HTTP
```
pandas-cli messages read <channel_id> <thing_auth_token>
```
