// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// Allows you to add filters when you use the search function in Secrets Manager.
// For more information, see Find secrets in Secrets Manager (https://docs.aws.amazon.com/secretsmanager/latest/userguide/manage_search-secret.html)
// .
type Filter struct {

	// The following are keys you can use:
	//   - description: Prefix match, not case-sensitive.
	//   - name: Prefix match, case-sensitive.
	//   - tag-key: Prefix match, case-sensitive.
	//   - tag-value: Prefix match, case-sensitive.
	//   - primary-region: Prefix match, case-sensitive.
	//   - owning-service: Prefix match, case-sensitive.
	//   - all: Breaks the filter value string into words and then searches all
	//   attributes for matches. Not case-sensitive.
	Key FilterNameStringType

	// The keyword to filter for. You can prefix your search value with an exclamation
	// mark ( ! ) in order to perform negation filters.
	Values []string

	noSmithyDocumentSerde
}

// A custom type that specifies a Region and the KmsKeyId for a replica secret.
type ReplicaRegionType struct {

	// The ARN, key ID, or alias of the KMS key to encrypt the secret. If you don't
	// include this field, Secrets Manager uses aws/secretsmanager .
	KmsKeyId *string

	// A Region code. For a list of Region codes, see Name and code of Regions (https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints)
	// .
	Region *string

	noSmithyDocumentSerde
}

// A replication object consisting of a RegionReplicationStatus object and
// includes a Region, KMSKeyId, status, and status message.
type ReplicationStatusType struct {

	// Can be an ARN , Key ID , or Alias .
	KmsKeyId *string

	// The date that the secret was last accessed in the Region. This field is omitted
	// if the secret has never been retrieved in the Region.
	LastAccessedDate *time.Time

	// The Region where replication occurs.
	Region *string

	// The status can be InProgress , Failed , or InSync .
	Status StatusType

	// Status message such as "Secret with this name already exists in this region".
	StatusMessage *string

	noSmithyDocumentSerde
}

// A structure that defines the rotation configuration for the secret.
type RotationRulesType struct {

	// The number of days between rotations of the secret. You can use this value to
	// check that your secret meets your compliance guidelines for how often secrets
	// must be rotated. If you use this field to set the rotation schedule, Secrets
	// Manager calculates the next rotation date based on the previous rotation.
	// Manually updating the secret value by calling PutSecretValue or UpdateSecret is
	// considered a valid rotation. In DescribeSecret and ListSecrets , this value is
	// calculated from the rotation schedule after every successful rotation. In
	// RotateSecret , you can set the rotation schedule in RotationRules with
	// AutomaticallyAfterDays or ScheduleExpression , but not both. To set a rotation
	// schedule in hours, use ScheduleExpression .
	AutomaticallyAfterDays *int64

	// The length of the rotation window in hours, for example 3h for a three hour
	// window. Secrets Manager rotates your secret at any time during this window. The
	// window must not extend into the next rotation window or the next UTC day. The
	// window starts according to the ScheduleExpression . If you don't specify a
	// Duration , for a ScheduleExpression in hours, the window automatically closes
	// after one hour. For a ScheduleExpression in days, the window automatically
	// closes at the end of the UTC day. For more information, including examples, see
	// Schedule expressions in Secrets Manager rotation (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotate-secrets_schedule.html)
	// in the Secrets Manager Users Guide.
	Duration *string

	// A cron() or rate() expression that defines the schedule for rotating your
	// secret. Secrets Manager rotation schedules use UTC time zone. Secrets Manager
	// rotates your secret any time during a rotation window. Secrets Manager rate()
	// expressions represent the interval in hours or days that you want to rotate your
	// secret, for example rate(12 hours) or rate(10 days) . You can rotate a secret as
	// often as every four hours. If you use a rate() expression, the rotation window
	// starts at midnight. For a rate in hours, the default rotation window closes
	// after one hour. For a rate in days, the default rotation window closes at the
	// end of the day. You can set the Duration to change the rotation window. The
	// rotation window must not extend into the next UTC day or into the next rotation
	// window. You can use a cron() expression to create a rotation schedule that is
	// more detailed than a rotation interval. For more information, including
	// examples, see Schedule expressions in Secrets Manager rotation (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotate-secrets_schedule.html)
	// in the Secrets Manager Users Guide. For a cron expression that represents a
	// schedule in hours, the default rotation window closes after one hour. For a cron
	// expression that represents a schedule in days, the default rotation window
	// closes at the end of the day. You can set the Duration to change the rotation
	// window. The rotation window must not extend into the next UTC day or into the
	// next rotation window.
	ScheduleExpression *string

	noSmithyDocumentSerde
}

// A structure that contains the details about a secret. It does not include the
// encrypted SecretString and SecretBinary values. To get those values, use
// GetSecretValue (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html)
// .
type SecretListEntry struct {

	// The Amazon Resource Name (ARN) of the secret.
	ARN *string

	// The date and time when a secret was created.
	CreatedDate *time.Time

	// The date and time the deletion of the secret occurred. Not present on active
	// secrets. The secret can be recovered until the number of days in the recovery
	// window has passed, as specified in the RecoveryWindowInDays parameter of the
	// DeleteSecret (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_DeleteSecret.html)
	// operation.
	DeletedDate *time.Time

	// The user-provided description of the secret.
	Description *string

	// The ARN of the KMS key that Secrets Manager uses to encrypt the secret value.
	// If the secret is encrypted with the Amazon Web Services managed key
	// aws/secretsmanager , this field is omitted.
	KmsKeyId *string

	// The date that the secret was last accessed in the Region. This field is omitted
	// if the secret has never been retrieved in the Region.
	LastAccessedDate *time.Time

	// The last date and time that this secret was modified in any way.
	LastChangedDate *time.Time

	// The most recent date and time that the Secrets Manager rotation process was
	// successfully completed. This value is null if the secret hasn't ever rotated.
	LastRotatedDate *time.Time

	// The friendly name of the secret. You can use forward slashes in the name to
	// represent a path hierarchy. For example, /prod/databases/dbserver1 could
	// represent the secret for a server named dbserver1 in the folder databases in
	// the folder prod .
	Name *string

	// The next date and time that Secrets Manager will attempt to rotate the secret,
	// rounded to the nearest hour. This value is null if the secret is not set up for
	// rotation.
	NextRotationDate *time.Time

	// Returns the name of the service that created the secret.
	OwningService *string

	// The Region where Secrets Manager originated the secret.
	PrimaryRegion *string

	// Indicates whether automatic, scheduled rotation is enabled for this secret.
	RotationEnabled *bool

	// The ARN of an Amazon Web Services Lambda function invoked by Secrets Manager to
	// rotate and expire the secret either automatically per the schedule or manually
	// by a call to RotateSecret (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_RotateSecret.html)
	// .
	RotationLambdaARN *string

	// A structure that defines the rotation configuration for the secret.
	RotationRules *RotationRulesType

	// A list of all of the currently assigned SecretVersionStage staging labels and
	// the SecretVersionId attached to each one. Staging labels are used to keep track
	// of the different versions during the rotation process. A version that does not
	// have any SecretVersionStage is considered deprecated and subject to deletion.
	// Such versions are not included in this list.
	SecretVersionsToStages map[string][]string

	// The list of user-defined tags associated with the secret. To add tags to a
	// secret, use TagResource (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_TagResource.html)
	// . To remove tags, use UntagResource (https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_UntagResource.html)
	// .
	Tags []Tag

	noSmithyDocumentSerde
}

// A structure that contains information about one version of a secret.
type SecretVersionsListEntry struct {

	// The date and time this version of the secret was created.
	CreatedDate *time.Time

	// The KMS keys used to encrypt the secret version.
	KmsKeyIds []string

	// The date that this version of the secret was last accessed. Note that the
	// resolution of this field is at the date level and does not include the time.
	LastAccessedDate *time.Time

	// The unique version identifier of this version of the secret.
	VersionId *string

	// An array of staging labels that are currently associated with this version of
	// the secret.
	VersionStages []string

	noSmithyDocumentSerde
}

// A structure that contains information about a tag.
type Tag struct {

	// The key identifier, or name, of the tag.
	Key *string

	// The string value associated with the key of the tag.
	Value *string

	noSmithyDocumentSerde
}

// Displays errors that occurred during validation of the resource policy.
type ValidationErrorsEntry struct {

	// Checks the name of the policy.
	CheckName *string

	// Displays error messages if validation encounters problems during validation of
	// the resource policy.
	ErrorMessage *string

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
