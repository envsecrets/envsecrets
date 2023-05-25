// Code generated by smithy-go-codegen DO NOT EDIT.

package secretsmanager

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Configures and starts the asynchronous process of rotating the secret. For
// information about rotation, see Rotate secrets (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets.html)
// in the Secrets Manager User Guide. If you include the configuration parameters,
// the operation sets the values for the secret and then immediately starts a
// rotation. If you don't include the configuration parameters, the operation
// starts a rotation with the values already stored in the secret. When rotation is
// successful, the AWSPENDING staging label might be attached to the same version
// as the AWSCURRENT version, or it might not be attached to any version. If the
// AWSPENDING staging label is present but not attached to the same version as
// AWSCURRENT , then any later invocation of RotateSecret assumes that a previous
// rotation request is still in progress and returns an error. When rotation is
// unsuccessful, the AWSPENDING staging label might be attached to an empty secret
// version. For more information, see Troubleshoot rotation (https://docs.aws.amazon.com/secretsmanager/latest/userguide/troubleshoot_rotation.html)
// in the Secrets Manager User Guide. Secrets Manager generates a CloudTrail log
// entry when you call this action. Do not include sensitive information in request
// parameters because it might be logged. For more information, see Logging
// Secrets Manager events with CloudTrail (https://docs.aws.amazon.com/secretsmanager/latest/userguide/retrieve-ct-entries.html)
// . Required permissions: secretsmanager:RotateSecret . For more information, see
// IAM policy actions for Secrets Manager (https://docs.aws.amazon.com/secretsmanager/latest/userguide/reference_iam-permissions.html#reference_iam-permissions_actions)
// and Authentication and access control in Secrets Manager (https://docs.aws.amazon.com/secretsmanager/latest/userguide/auth-and-access.html)
// . You also need lambda:InvokeFunction permissions on the rotation function. For
// more information, see Permissions for rotation (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets-required-permissions-function.html)
// .
func (c *Client) RotateSecret(ctx context.Context, params *RotateSecretInput, optFns ...func(*Options)) (*RotateSecretOutput, error) {
	if params == nil {
		params = &RotateSecretInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "RotateSecret", params, optFns, c.addOperationRotateSecretMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*RotateSecretOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type RotateSecretInput struct {

	// The ARN or name of the secret to rotate. For an ARN, we recommend that you
	// specify a complete ARN rather than a partial ARN. See Finding a secret from a
	// partial ARN (https://docs.aws.amazon.com/secretsmanager/latest/userguide/troubleshoot.html#ARN_secretnamehyphen)
	// .
	//
	// This member is required.
	SecretId *string

	// A unique identifier for the new version of the secret that helps ensure
	// idempotency. Secrets Manager uses this value to prevent the accidental creation
	// of duplicate versions if there are failures and retries during rotation. This
	// value becomes the VersionId of the new version. If you use the Amazon Web
	// Services CLI or one of the Amazon Web Services SDK to call this operation, then
	// you can leave this parameter empty. The CLI or SDK generates a random UUID for
	// you and includes that in the request for this parameter. If you don't use the
	// SDK and instead generate a raw HTTP request to the Secrets Manager service
	// endpoint, then you must generate a ClientRequestToken yourself for new versions
	// and include that value in the request. You only need to specify this value if
	// you implement your own retry logic and you want to ensure that Secrets Manager
	// doesn't attempt to create a secret version twice. We recommend that you generate
	// a UUID-type (https://wikipedia.org/wiki/Universally_unique_identifier) value to
	// ensure uniqueness within the specified secret.
	ClientRequestToken *string

	// Specifies whether to rotate the secret immediately or wait until the next
	// scheduled rotation window. The rotation schedule is defined in
	// RotateSecretRequest$RotationRules . For secrets that use a Lambda rotation
	// function to rotate, if you don't immediately rotate the secret, Secrets Manager
	// tests the rotation configuration by running the testSecret step (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotate-secrets_how.html)
	// of the Lambda rotation function. The test creates an AWSPENDING version of the
	// secret and then removes it. By default, Secrets Manager rotates the secret
	// immediately.
	RotateImmediately *bool

	// For secrets that use a Lambda rotation function to rotate, the ARN of the
	// Lambda rotation function. For secrets that use managed rotation, omit this
	// field. For more information, see Managed rotation (https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotate-secrets_managed.html)
	// in the Secrets Manager User Guide.
	RotationLambdaARN *string

	// A structure that defines the rotation configuration for this secret.
	RotationRules *types.RotationRulesType

	noSmithyDocumentSerde
}

type RotateSecretOutput struct {

	// The ARN of the secret.
	ARN *string

	// The name of the secret.
	Name *string

	// The ID of the new version of the secret.
	VersionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationRotateSecretMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpRotateSecret{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpRotateSecret{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addIdempotencyToken_opRotateSecretMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpRotateSecretValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opRotateSecret(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

type idempotencyToken_initializeOpRotateSecret struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpRotateSecret) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpRotateSecret) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*RotateSecretInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *RotateSecretInput ")
	}

	if input.ClientRequestToken == nil {
		t, err := m.tokenProvider.GetIdempotencyToken()
		if err != nil {
			return out, metadata, err
		}
		input.ClientRequestToken = &t
	}
	return next.HandleInitialize(ctx, in)
}
func addIdempotencyToken_opRotateSecretMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpRotateSecret{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opRotateSecret(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "secretsmanager",
		OperationName: "RotateSecret",
	}
}
