package databrew

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Recipe_RecipeParameters AWS CloudFormation Resource (AWS::DataBrew::Recipe.RecipeParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html
type Recipe_RecipeParameters struct {

	// AggregateFunction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-aggregatefunction
	AggregateFunction *types.Value `json:"AggregateFunction,omitempty"`

	// Base AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-base
	Base *types.Value `json:"Base,omitempty"`

	// CaseStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-casestatement
	CaseStatement *types.Value `json:"CaseStatement,omitempty"`

	// CategoryMap AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-categorymap
	CategoryMap *types.Value `json:"CategoryMap,omitempty"`

	// CharsToRemove AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-charstoremove
	CharsToRemove *types.Value `json:"CharsToRemove,omitempty"`

	// CollapseConsecutiveWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-collapseconsecutivewhitespace
	CollapseConsecutiveWhitespace *types.Value `json:"CollapseConsecutiveWhitespace,omitempty"`

	// ColumnDataType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-columndatatype
	ColumnDataType *types.Value `json:"ColumnDataType,omitempty"`

	// ColumnRange AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-columnrange
	ColumnRange *types.Value `json:"ColumnRange,omitempty"`

	// Count AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-count
	Count *types.Value `json:"Count,omitempty"`

	// CustomCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customcharacters
	CustomCharacters *types.Value `json:"CustomCharacters,omitempty"`

	// CustomStopWords AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customstopwords
	CustomStopWords *types.Value `json:"CustomStopWords,omitempty"`

	// CustomValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customvalue
	CustomValue *types.Value `json:"CustomValue,omitempty"`

	// DatasetsColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datasetscolumns
	DatasetsColumns *types.Value `json:"DatasetsColumns,omitempty"`

	// DateAddValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-dateaddvalue
	DateAddValue *types.Value `json:"DateAddValue,omitempty"`

	// DateTimeFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datetimeformat
	DateTimeFormat *types.Value `json:"DateTimeFormat,omitempty"`

	// DateTimeParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datetimeparameters
	DateTimeParameters *types.Value `json:"DateTimeParameters,omitempty"`

	// DeleteOtherRows AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-deleteotherrows
	DeleteOtherRows *types.Value `json:"DeleteOtherRows,omitempty"`

	// Delimiter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-delimiter
	Delimiter *types.Value `json:"Delimiter,omitempty"`

	// EndPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endpattern
	EndPattern *types.Value `json:"EndPattern,omitempty"`

	// EndPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endposition
	EndPosition *types.Value `json:"EndPosition,omitempty"`

	// EndValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endvalue
	EndValue *types.Value `json:"EndValue,omitempty"`

	// ExpandContractions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-expandcontractions
	ExpandContractions *types.Value `json:"ExpandContractions,omitempty"`

	// Exponent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-exponent
	Exponent *types.Value `json:"Exponent,omitempty"`

	// FalseString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-falsestring
	FalseString *types.Value `json:"FalseString,omitempty"`

	// GroupByAggFunctionOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-groupbyaggfunctionoptions
	GroupByAggFunctionOptions *types.Value `json:"GroupByAggFunctionOptions,omitempty"`

	// GroupByColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-groupbycolumns
	GroupByColumns *types.Value `json:"GroupByColumns,omitempty"`

	// HiddenColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-hiddencolumns
	HiddenColumns *types.Value `json:"HiddenColumns,omitempty"`

	// IgnoreCase AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-ignorecase
	IgnoreCase *types.Value `json:"IgnoreCase,omitempty"`

	// IncludeInSplit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-includeinsplit
	IncludeInSplit *types.Value `json:"IncludeInSplit,omitempty"`

	// Input AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-input
	Input interface{} `json:"Input,omitempty"`

	// Interval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-interval
	Interval *types.Value `json:"Interval,omitempty"`

	// IsText AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-istext
	IsText *types.Value `json:"IsText,omitempty"`

	// JoinKeys AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-joinkeys
	JoinKeys *types.Value `json:"JoinKeys,omitempty"`

	// JoinType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-jointype
	JoinType *types.Value `json:"JoinType,omitempty"`

	// LeftColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-leftcolumns
	LeftColumns *types.Value `json:"LeftColumns,omitempty"`

	// Limit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-limit
	Limit *types.Value `json:"Limit,omitempty"`

	// LowerBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-lowerbound
	LowerBound *types.Value `json:"LowerBound,omitempty"`

	// MapType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-maptype
	MapType *types.Value `json:"MapType,omitempty"`

	// ModeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-modetype
	ModeType *types.Value `json:"ModeType,omitempty"`

	// MultiLine AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-multiline
	MultiLine *types.Value `json:"MultiLine,omitempty"`

	// NumRows AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrows
	NumRows *types.Value `json:"NumRows,omitempty"`

	// NumRowsAfter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrowsafter
	NumRowsAfter *types.Value `json:"NumRowsAfter,omitempty"`

	// NumRowsBefore AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrowsbefore
	NumRowsBefore *types.Value `json:"NumRowsBefore,omitempty"`

	// OrderByColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-orderbycolumn
	OrderByColumn *types.Value `json:"OrderByColumn,omitempty"`

	// OrderByColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-orderbycolumns
	OrderByColumns *types.Value `json:"OrderByColumns,omitempty"`

	// Other AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-other
	Other *types.Value `json:"Other,omitempty"`

	// Pattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-pattern
	Pattern *types.Value `json:"Pattern,omitempty"`

	// PatternOption1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoption1
	PatternOption1 *types.Value `json:"PatternOption1,omitempty"`

	// PatternOption2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoption2
	PatternOption2 *types.Value `json:"PatternOption2,omitempty"`

	// PatternOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoptions
	PatternOptions *types.Value `json:"PatternOptions,omitempty"`

	// Period AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-period
	Period *types.Value `json:"Period,omitempty"`

	// Position AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-position
	Position *types.Value `json:"Position,omitempty"`

	// RemoveAllPunctuation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallpunctuation
	RemoveAllPunctuation *types.Value `json:"RemoveAllPunctuation,omitempty"`

	// RemoveAllQuotes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallquotes
	RemoveAllQuotes *types.Value `json:"RemoveAllQuotes,omitempty"`

	// RemoveAllWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallwhitespace
	RemoveAllWhitespace *types.Value `json:"RemoveAllWhitespace,omitempty"`

	// RemoveCustomCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removecustomcharacters
	RemoveCustomCharacters *types.Value `json:"RemoveCustomCharacters,omitempty"`

	// RemoveCustomValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removecustomvalue
	RemoveCustomValue *types.Value `json:"RemoveCustomValue,omitempty"`

	// RemoveLeadingAndTrailingPunctuation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingpunctuation
	RemoveLeadingAndTrailingPunctuation *types.Value `json:"RemoveLeadingAndTrailingPunctuation,omitempty"`

	// RemoveLeadingAndTrailingQuotes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingquotes
	RemoveLeadingAndTrailingQuotes *types.Value `json:"RemoveLeadingAndTrailingQuotes,omitempty"`

	// RemoveLeadingAndTrailingWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingwhitespace
	RemoveLeadingAndTrailingWhitespace *types.Value `json:"RemoveLeadingAndTrailingWhitespace,omitempty"`

	// RemoveLetters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeletters
	RemoveLetters *types.Value `json:"RemoveLetters,omitempty"`

	// RemoveNumbers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removenumbers
	RemoveNumbers *types.Value `json:"RemoveNumbers,omitempty"`

	// RemoveSourceColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removesourcecolumn
	RemoveSourceColumn *types.Value `json:"RemoveSourceColumn,omitempty"`

	// RemoveSpecialCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removespecialcharacters
	RemoveSpecialCharacters *types.Value `json:"RemoveSpecialCharacters,omitempty"`

	// RightColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-rightcolumns
	RightColumns *types.Value `json:"RightColumns,omitempty"`

	// SampleSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-samplesize
	SampleSize *types.Value `json:"SampleSize,omitempty"`

	// SampleType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sampletype
	SampleType *types.Value `json:"SampleType,omitempty"`

	// SecondInput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-secondinput
	SecondInput *types.Value `json:"SecondInput,omitempty"`

	// SecondaryInputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-secondaryinputs
	SecondaryInputs []Recipe_SecondaryInput `json:"SecondaryInputs,omitempty"`

	// SheetIndexes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sheetindexes
	SheetIndexes *types.Value `json:"SheetIndexes,omitempty"`

	// SheetNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sheetnames
	SheetNames *types.Value `json:"SheetNames,omitempty"`

	// SourceColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn
	SourceColumn *types.Value `json:"SourceColumn,omitempty"`

	// SourceColumn1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn1
	SourceColumn1 *types.Value `json:"SourceColumn1,omitempty"`

	// SourceColumn2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn2
	SourceColumn2 *types.Value `json:"SourceColumn2,omitempty"`

	// SourceColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumns
	SourceColumns *types.Value `json:"SourceColumns,omitempty"`

	// StartColumnIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startcolumnindex
	StartColumnIndex *types.Value `json:"StartColumnIndex,omitempty"`

	// StartPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startpattern
	StartPattern *types.Value `json:"StartPattern,omitempty"`

	// StartPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startposition
	StartPosition *types.Value `json:"StartPosition,omitempty"`

	// StartValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startvalue
	StartValue *types.Value `json:"StartValue,omitempty"`

	// StemmingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stemmingmode
	StemmingMode *types.Value `json:"StemmingMode,omitempty"`

	// StepCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stepcount
	StepCount *types.Value `json:"StepCount,omitempty"`

	// StepIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stepindex
	StepIndex *types.Value `json:"StepIndex,omitempty"`

	// StopWordsMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stopwordsmode
	StopWordsMode *types.Value `json:"StopWordsMode,omitempty"`

	// Strategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-strategy
	Strategy *types.Value `json:"Strategy,omitempty"`

	// TargetColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetcolumn
	TargetColumn *types.Value `json:"TargetColumn,omitempty"`

	// TargetColumnNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetcolumnnames
	TargetColumnNames *types.Value `json:"TargetColumnNames,omitempty"`

	// TargetDateFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetdateformat
	TargetDateFormat *types.Value `json:"TargetDateFormat,omitempty"`

	// TargetIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetindex
	TargetIndex *types.Value `json:"TargetIndex,omitempty"`

	// TimeZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-timezone
	TimeZone *types.Value `json:"TimeZone,omitempty"`

	// TokenizerPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-tokenizerpattern
	TokenizerPattern *types.Value `json:"TokenizerPattern,omitempty"`

	// TrueString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-truestring
	TrueString *types.Value `json:"TrueString,omitempty"`

	// UdfLang AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-udflang
	UdfLang *types.Value `json:"UdfLang,omitempty"`

	// Units AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-units
	Units *types.Value `json:"Units,omitempty"`

	// UnpivotColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-unpivotcolumn
	UnpivotColumn *types.Value `json:"UnpivotColumn,omitempty"`

	// UpperBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-upperbound
	UpperBound *types.Value `json:"UpperBound,omitempty"`

	// UseNewDataFrame AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-usenewdataframe
	UseNewDataFrame *types.Value `json:"UseNewDataFrame,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value
	Value *types.Value `json:"Value,omitempty"`

	// Value1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value1
	Value1 *types.Value `json:"Value1,omitempty"`

	// Value2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value2
	Value2 *types.Value `json:"Value2,omitempty"`

	// ValueColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-valuecolumn
	ValueColumn *types.Value `json:"ValueColumn,omitempty"`

	// ViewFrame AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-viewframe
	ViewFrame *types.Value `json:"ViewFrame,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Recipe_RecipeParameters) AWSCloudFormationType() string {
	return "AWS::DataBrew::Recipe.RecipeParameters"
}
