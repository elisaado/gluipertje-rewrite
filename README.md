# gluipertje-rewrite
The rewritten backend for Gluipertje

## Usage and docs
All routes are prefixed with /api.

### Types

##### User

| Field | Type | Public |
| ----- | ---- | ------ |
| ID | Int (unique) | Yes |
| Created at | Time object | Yes |
| Updated at | Time object | Yes |
| Deleted at | Time object | Yes |
| Nickname | String (length < 50) | Yes |
| Username | String (unique, length < 50) | Yes |
| Password | String | No |
| Token | String (length = 60) | No |

##### Message
| Field | Type | Public |
| ----- | ---- | ------ |
| ID | Int (unique) | Yes |
| Created at | Time object | Yes |
| Updated at | Time object | Yes |
| Deleted at | Time object | Yes |
| Type | String | Yes |
| Text | String (length <= 500) | Yes |
| From ID | Int (User ID) | Yes |
| From | User | Yes|

**Note: Type is either `text` or `image`. A message with type `image` may still have text (e.g. a caption)**

### Methods
POST request have to be JSON encoded (and must include a `Content-Type: application/json` header).

| HTTP verb | Method | Params | Returns | Action | Protected |
| --------- | ------ | ------ | ------- | ------ | --------- |
| GET | /users | None | List of users | Get all users | No |
| POST | /users | Nickname, Username and Password from the User type | Newly created user | Create a new User | No |
| GET | /user/:id | User ID (int) | User | Get specific user by ID | No |
| GET | /user/:username | User Username (string) | User | Get specific user by Username | No |
| GET | /messages | None | List of messages | Get all messages | No |
| POST | /messages | Body from Message type, Token from User type | Newly created message | Send/Create a new Message | Yes |
| GET | /messages/:limit | Limit (int) | List of messages | Get all messages limited by limit | No |
| GET | /message/:id | Message ID (int) | Message | Get specific message by ID | No |
| POST | /token | Username and Password from the User type | New Token (string) | Revoke the current Token | Yes |
| GET | /:token/me | Token from User type | Current User | Get info about the current User | Yes |

**Note: as you may have noticed, /user/:id and /user/:username are the same request. If the provided parameter consists only of numbers, it is interpreted as an ID, otherwise as a username (123 = ID, 123a = username). This is OK because usernames may not only exist of numbers.**