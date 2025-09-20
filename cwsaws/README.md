# cwsaws - AWS Services Utility Library

The cwsaws library provides a simplified interface to work with various AWS services in Go applications. It offers a consistent pattern for initializing service clients, connection pooling, and higher-level abstractions over the AWS SDK for Go v2.

## Supported AWS Services

- **DynamoDB**: CRUD operations, table management, and batch processing
- **S3**: Object storage with local caching support
- **SQS**: Message queue operations and Lambda event processing
- **SNS**: Notification and SMS delivery
- **SES**: Email delivery with template support
- **CloudWatch**: Logging and monitoring
- **STS**: Identity and access management

## Environment Variables

| Variable                  | Module | Type   | Description                                                 |
| ------------------------- | ------ | ------ | ----------------------------------------------------------- |
| S3CacheTTL                | cwsaws | int    | S3 Object local cache time to live in minutes (default: 10) |
| S3VersionCheck            | cwsaws | int    | Time in seconds between S3 version checks (default: 30)     |
| CLOUDWATCHLOG_LOG_GROUP   | cwsaws | string | AWS CloudWatch log group name                               |
| Local_DynamoDB_AWS_ID     | cwsaws | string | AWS ID for local DynamoDB connections                       |
| Local_DynamoDB_AWS_Secret | cwsaws | string | AWS Secret for local DynamoDB connections                   |
| Local_DynamoDB_URL        | cwsaws | string | URL for local DynamoDB endpoint                             |
| Local_DynamoDB_REGION     | cwsaws | string | Region for local DynamoDB connections                       |

## Installation

```bash
go get github.com/codeworks-tw/cwsutil/cwsaws
```

## Usage Examples

### DynamoDB

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Create a DynamoDB table proxy
ctx := context.Background()
tableProxy := cwsaws.GetDynamoDBTableProxy[map[string]any]("my-table", ctx)

// Get item from table
key := map[string]types.AttributeValue{
    "id": &types.AttributeValueMemberS{Value: "123"},
}
input := &dynamodb.GetItemInput{
    Key: key,
}
item, err := tableProxy.ProxyGetItem(input)
```

### S3

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Create an S3 proxy for a specific bucket
ctx := context.Background()
s3Proxy := cwsaws.GetS3Proxy(ctx, "my-bucket")

// Check if object exists
exists := s3Proxy.ProxyObjectExists("path/to/object.json")

// Get object with caching
jsonObject, err := s3Proxy.ProxyGetObject("path/to/object.json", func(content []byte) (any, error) {
    var data map[string]any
    err := json.Unmarshal(content, &data)
    return data, err
})
```

### SQS

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Create an SQS proxy and process messages
ctx := context.Background()
sqsProxy := cwsaws.GetSqsProxy(ctx, func(ctx context.Context, msg types.Message) error {
    // Process message
    fmt.Println("Message body:", *msg.Body)
    return nil
})

// Initialize by queue name
err := sqsProxy.ProxyInitializeByName("my-queue")

// Process next batch of messages
err = sqsProxy.ProxyProcessNextMessages(10)
```

### CloudWatch Logging

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Create a CloudWatch logs proxy
ctx := context.Background()
cwProxy := cwsaws.GetCloudWatchLogProxy(ctx)

// Send a log message
err := cwProxy.SendMessage("My log message")
```

### SES (Simple Email Service)

```go path=null start=null
import (
    "github.com/codeworks-tw/cwsutil/cwsaws"
    "net/mail"
)

// Create an SES proxy
ctx := context.Background()
sesProxy := cwsaws.GetSesProxy(ctx)

// Send an email
email := &cwsaws.SESProxyEmailInput{
    From: mail.Address{Name: "Sender", Address: "sender@example.com"},
    To: []mail.Address{{Name: "Recipient", Address: "recipient@example.com"}},
    Subject: "Test Email",
    Body: "<h1>Hello World</h1>",
    // IsText: false, // Default is HTML
}
result, err := sesProxy.ProxySendEmail(email)
```

### SNS (Simple Notification Service)

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Create an SNS proxy
ctx := context.Background()
snsProxy := cwsaws.GetSnsProxy(ctx)

// Send SMS notification
result, err := snsProxy.ProxySendPhoneNotification("+1234567890", "Hello World from SNS")

// Send notification using template
templateInput := &cwsaws.SNSProxyTemplateNotificationInput{
    TemplateName: "my-template",
    Phone:        "+1234567890",
    Params:       []any{"John", "Doe"},
}
result, err = snsProxy.ProxySendTemplateNotification(templateInput)
```

## Generic Repository Pattern

The cwsaws library provides a generic Repository pattern for DynamoDB operations:

```go path=null start=null
import "github.com/codeworks-tw/cwsutil/cwsaws"

// Define primary key structure
type UserPKey struct {
    ID string `json:"id"`
}

// Create repository
repo := &cwsaws.Repository[UserPKey]{
    TableName: "users",
}

// Get item
ctx := context.Background()
pKey := UserPKey{ID: "123"}
user, err := repo.Get(ctx, pKey)

// Update item with expression
updateExpr := expression.Set(expression.Name("status"), expression.Value("active"))
expr, _ := expression.NewBuilder().WithUpdate(updateExpr).Build()
updatedUser, err := repo.Merge(ctx, pKey, expr)

// Query items
queryExpr := expression.Key("id").Equal(expression.Value("123"))
keyExpr, _ := expression.NewBuilder().WithKeyCondition(queryExpr).Build()
items, err := repo.Query(ctx, "my-index", keyExpr)

// Delete item
deletedUser, err := repo.Delete(ctx, pKey)
```

## Release History

* 0.3.7 - Sep. 20, 2025 - Added CloudWatch logging support
* 0.3.6 - Jun. 15, 2025
* 0.3.5 - May. 25, 2025
* 0.1.0 - Apr. 11, 2024 - Initial release

## License

Copyright (c) 2024 - Present Codeworks TW Ltd.
