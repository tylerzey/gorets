package metadata

// MValidationExpression ...
type MValidationExpression struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource string   `xml:"Resource,attr"`

	ValidationExpression []ValidationExpression `xml:"ValidationExpression"`
}

// ValidationExpression ...
type ValidationExpression struct {
	MetadataEntryID          string `xml:"MetadataEntryID"`
	ValidationExpressionID   string `xml:"ValidationExpressionID"`
	ValidationExpressionType string `xml:"ValidationExpressionType"`
	Message                  string `xml:"Message"`
	IsCaseSensitive          *bool  `xml:"IsCaseSensitive"`
}
