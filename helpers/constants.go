package helpers

const Name = "cbgo"

const (
	Prefix = "_connector:" + Name + ":"

	DefaultScopeName                            = "_default"
	DefaultCollectionName                       = "_default"
	MetadataTypeCouchbase                       = "couchbase"
	MetadataTypeFile                            = "file"
	CouchbaseMetadataBucketConfig               = "bucket"
	CouchbaseMetadataScopeConfig                = "scope"
	CouchbaseMetadataCollectionConfig           = "collection"
	CouchbaseMetadataConnectionBufferSizeConfig = "connectionBufferSize"
	CouchbaseMetadataConnectionTimeoutConfig    = "connectionTimeout"
	FileMetadataFileNameConfig                  = "fileName"
	MembershipChangedBusEventName               = "membershipChanged"
	PersistSeqNoChangedBusEventName             = "persistSeqNoChanged"
)
