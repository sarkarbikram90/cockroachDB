// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists the server certificates stored in IAM that have the specified path prefix.
// If none exist, the operation returns an empty list. You can paginate the results
// using the MaxItems and Marker parameters. For more information about working
// with server certificates, see Working with server certificates
// (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_server-certs.html)
// in the IAM User Guide. This topic also includes a list of Amazon Web Services
// services that can use the server certificates that you manage with IAM. IAM
// resource-listing operations return a subset of the available attributes for the
// resource. For example, this operation does not return tags, even though they are
// an attribute of the returned object. To view all of the information for a
// servercertificate, see GetServerCertificate.
func (c *Client) ListServerCertificates(ctx context.Context, params *ListServerCertificatesInput, optFns ...func(*Options)) (*ListServerCertificatesOutput, error) {
	if params == nil {
		params = &ListServerCertificatesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListServerCertificates", params, optFns, c.addOperationListServerCertificatesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListServerCertificatesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListServerCertificatesInput struct {

	// Use this parameter only when paginating results and only after you receive a
	// response indicating that the results are truncated. Set it to the value of the
	// Marker element in the response that you received to indicate where the next call
	// should start.
	Marker *string

	// Use this only when paginating results to indicate the maximum number of items
	// you want in the response. If additional items exist beyond the maximum you
	// specify, the IsTruncated response element is true. If you do not include this
	// parameter, the number of items defaults to 100. Note that IAM might return fewer
	// results, even when there are more results available. In that case, the
	// IsTruncated response element returns true, and Marker contains a value to
	// include in the subsequent call that tells the service where to continue from.
	MaxItems *int32

	// The path prefix for filtering the results. For example: /company/servercerts
	// would get all server certificates for which the path starts with
	// /company/servercerts. This parameter is optional. If it is not included, it
	// defaults to a slash (/), listing all server certificates. This parameter allows
	// (through its regex pattern (http://wikipedia.org/wiki/regex)) a string of
	// characters consisting of either a forward slash (/) by itself or a string that
	// must begin and end with forward slashes. In addition, it can contain any ASCII
	// character from the ! (\u0021) through the DEL character (\u007F), including most
	// punctuation characters, digits, and upper and lowercased letters.
	PathPrefix *string

	noSmithyDocumentSerde
}

// Contains the response to a successful ListServerCertificates request.
type ListServerCertificatesOutput struct {

	// A list of server certificates.
	//
	// This member is required.
	ServerCertificateMetadataList []types.ServerCertificateMetadata

	// A flag that indicates whether there are more items to return. If your results
	// were truncated, you can make a subsequent pagination request using the Marker
	// request parameter to retrieve more items. Note that IAM might return fewer than
	// the MaxItems number of results even when there are more results available. We
	// recommend that you check IsTruncated after every call to ensure that you receive
	// all your results.
	IsTruncated bool

	// When IsTruncated is true, this element is present and contains the value to use
	// for the Marker parameter in a subsequent pagination request.
	Marker *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListServerCertificatesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpListServerCertificates{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpListServerCertificates{}, middleware.After)
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListServerCertificates(options.Region), middleware.Before); err != nil {
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

// ListServerCertificatesAPIClient is a client that implements the
// ListServerCertificates operation.
type ListServerCertificatesAPIClient interface {
	ListServerCertificates(context.Context, *ListServerCertificatesInput, ...func(*Options)) (*ListServerCertificatesOutput, error)
}

var _ ListServerCertificatesAPIClient = (*Client)(nil)

// ListServerCertificatesPaginatorOptions is the paginator options for
// ListServerCertificates
type ListServerCertificatesPaginatorOptions struct {
	// Use this only when paginating results to indicate the maximum number of items
	// you want in the response. If additional items exist beyond the maximum you
	// specify, the IsTruncated response element is true. If you do not include this
	// parameter, the number of items defaults to 100. Note that IAM might return fewer
	// results, even when there are more results available. In that case, the
	// IsTruncated response element returns true, and Marker contains a value to
	// include in the subsequent call that tells the service where to continue from.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListServerCertificatesPaginator is a paginator for ListServerCertificates
type ListServerCertificatesPaginator struct {
	options   ListServerCertificatesPaginatorOptions
	client    ListServerCertificatesAPIClient
	params    *ListServerCertificatesInput
	nextToken *string
	firstPage bool
}

// NewListServerCertificatesPaginator returns a new ListServerCertificatesPaginator
func NewListServerCertificatesPaginator(client ListServerCertificatesAPIClient, params *ListServerCertificatesInput, optFns ...func(*ListServerCertificatesPaginatorOptions)) *ListServerCertificatesPaginator {
	if params == nil {
		params = &ListServerCertificatesInput{}
	}

	options := ListServerCertificatesPaginatorOptions{}
	if params.MaxItems != nil {
		options.Limit = *params.MaxItems
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListServerCertificatesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListServerCertificatesPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next ListServerCertificates page.
func (p *ListServerCertificatesPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListServerCertificatesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.Marker = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxItems = limit

	result, err := p.client.ListServerCertificates(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.Marker

	if p.options.StopOnDuplicateToken && prevToken != nil && p.nextToken != nil && *prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListServerCertificates(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "iam",
		OperationName: "ListServerCertificates",
	}
}
