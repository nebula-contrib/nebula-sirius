/*
 *
 * Copyright (c) 2020 Elchin Gasimov. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 *
 * This file contents is copied from vesoft-inc/nebula-go library, make necessary changes
 */

package nebula_sirius

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/nebula"
	"github.com/nebula-contrib/nebula-sirius/nebula/graph"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

type ResultSet struct {
	resp            *graph.ExecutionResponse
	columnNames     []string
	colNameIndexMap map[string]int
	timezoneInfo    timezoneInfo
}

type Record struct {
	columnNames     *[]string
	_record         []*ValueWrapper
	colNameIndexMap *map[string]int
	timezoneInfo    timezoneInfo
}

type Node struct {
	vertex          *nebula.Vertex
	tags            []string // tag name
	tagNameIndexMap map[string]int
	timezoneInfo    timezoneInfo
}

type Relationship struct {
	edge         *nebula.Edge
	timezoneInfo timezoneInfo
}

type segment struct {
	startNode    *Node
	relationship *Relationship
	endNode      *Node
}

type PathWrapper struct {
	path             *nebula.Path
	nodeList         []*Node
	relationshipList []*Relationship
	segments         []segment
	timezoneInfo     timezoneInfo
}

type TimeWrapper struct {
	time         *nebula.Time
	timezoneInfo timezoneInfo
}

type DateWrapper struct {
	date *nebula.Date
}

type DateTimeWrapper struct {
	dateTime     *nebula.DateTime
	timezoneInfo timezoneInfo
}

type ErrorCode int64

const (
	ErrorCode_SUCCEEDED               ErrorCode = ErrorCode(nebula.ErrorCode_SUCCEEDED)
	ErrorCode_E_DISCONNECTED          ErrorCode = ErrorCode(nebula.ErrorCode_E_DISCONNECTED)
	ErrorCode_E_FAIL_TO_CONNECT       ErrorCode = ErrorCode(nebula.ErrorCode_E_FAIL_TO_CONNECT)
	ErrorCode_E_RPC_FAILURE           ErrorCode = ErrorCode(nebula.ErrorCode_E_RPC_FAILURE)
	ErrorCode_E_BAD_USERNAME_PASSWORD ErrorCode = ErrorCode(nebula.ErrorCode_E_BAD_USERNAME_PASSWORD)
	ErrorCode_E_SESSION_INVALID       ErrorCode = ErrorCode(nebula.ErrorCode_E_SESSION_INVALID)
	ErrorCode_E_SESSION_TIMEOUT       ErrorCode = ErrorCode(nebula.ErrorCode_E_SESSION_TIMEOUT)
	ErrorCode_E_SYNTAX_ERROR          ErrorCode = ErrorCode(nebula.ErrorCode_E_SYNTAX_ERROR)
	ErrorCode_E_EXECUTION_ERROR       ErrorCode = ErrorCode(nebula.ErrorCode_E_EXECUTION_ERROR)
	ErrorCode_E_STATEMENT_EMPTY       ErrorCode = ErrorCode(nebula.ErrorCode_E_STATEMENT_EMPTY)
	ErrorCode_E_USER_NOT_FOUND        ErrorCode = ErrorCode(nebula.ErrorCode_E_USER_NOT_FOUND)
	ErrorCode_E_BAD_PERMISSION        ErrorCode = ErrorCode(nebula.ErrorCode_E_BAD_PERMISSION)
	ErrorCode_E_SEMANTIC_ERROR        ErrorCode = ErrorCode(nebula.ErrorCode_E_SEMANTIC_ERROR)
	ErrorCode_E_PARTIAL_SUCCEEDED     ErrorCode = ErrorCode(nebula.ErrorCode_E_PARTIAL_SUCCEEDED)
)

func GenResultSet(resp *graph.ExecutionResponse) (*ResultSet, error) {
	var defaultTimezone timezoneInfo = timezoneInfo{0, []byte("UTC")}
	return genResultSet(resp, defaultTimezone)
}

func genResultSet(resp *graph.ExecutionResponse, timezoneInfo timezoneInfo) (*ResultSet, error) {
	var colNames []string
	var colNameIndexMap = make(map[string]int)

	if resp.Data == nil { // if resp.Data != nil then resp.Data.row and resp.Data.colNames wont be nil
		return &ResultSet{
			resp:            resp,
			columnNames:     colNames,
			colNameIndexMap: colNameIndexMap,
		}, nil
	}
	for i, name := range resp.Data.ColumnNames {
		colNames = append(colNames, string(name))
		colNameIndexMap[string(name)] = i
	}

	return &ResultSet{
		resp:            resp,
		columnNames:     colNames,
		colNameIndexMap: colNameIndexMap,
		timezoneInfo:    timezoneInfo,
	}, nil
}

func genValWraps(row *nebula.Row, timezoneInfo timezoneInfo) ([]*ValueWrapper, error) {
	if row == nil {
		return nil, fmt.Errorf("failed to generate valueWrapper: invalid row")
	}
	var valWraps []*ValueWrapper
	for _, val := range row.Values {
		if val == nil {
			return nil, fmt.Errorf("failed to generate valueWrapper: value is nil")
		}
		valWraps = append(valWraps, &ValueWrapper{val, timezoneInfo})
	}
	return valWraps, nil
}

func genNode(vertex *nebula.Vertex, timezoneInfo timezoneInfo) (*Node, error) {
	if vertex == nil {
		return nil, fmt.Errorf("failed to generate Node: invalid vertex")
	}
	var tags []string
	nameIndex := make(map[string]int)

	// Iterate through all tags of the vertex
	for i, tag := range vertex.GetTags() {
		name := string(tag.Name)
		// Get tags
		tags = append(tags, name)
		nameIndex[name] = i
	}

	return &Node{
		vertex:          vertex,
		tags:            tags,
		tagNameIndexMap: nameIndex,
		timezoneInfo:    timezoneInfo,
	}, nil
}

func genRelationship(edge *nebula.Edge, timezoneInfo timezoneInfo) (*Relationship, error) {
	if edge == nil {
		return nil, fmt.Errorf("failed to generate Relationship: invalid edge")
	}
	return &Relationship{
		edge:         edge,
		timezoneInfo: timezoneInfo,
	}, nil
}

func genPathWrapper(path *nebula.Path, timezoneInfo timezoneInfo) (*PathWrapper, error) {
	if path == nil {
		return nil, fmt.Errorf("failed to generate Path Wrapper: invalid path")
	}
	var (
		nodeList         []*Node
		relationshipList []*Relationship
		segList          []segment
		edge             *nebula.Edge
		segStartNode     *Node
		segEndNode       *Node
		segType          nebula.EdgeType
	)
	src, err := genNode(path.Src, timezoneInfo)
	if err != nil {
		return nil, err
	}
	nodeList = append(nodeList, src)

	for _, step := range path.Steps {
		dst, err := genNode(step.Dst, timezoneInfo)
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, dst)
		// determine direction
		stepType := step.Type
		if stepType > 0 {
			segStartNode = src
			segEndNode = dst
			segType = stepType
		} else {
			segStartNode = dst // switch with src
			segEndNode = src
			segType = -stepType
		}
		edge = &nebula.Edge{
			Src:     segStartNode.getRawID(),
			Dst:     segEndNode.getRawID(),
			Type:    segType,
			Name:    step.Name,
			Ranking: step.Ranking,
			Props:   step.Props,
		}
		relationship, err := genRelationship(edge, timezoneInfo)
		if err != nil {
			return nil, err
		}
		relationshipList = append(relationshipList, relationship)

		// Check segments
		if len(segList) > 0 {
			prevStart := segList[len(segList)-1].startNode.GetID()
			prevEnd := segList[len(segList)-1].endNode.GetID()
			nextStart := segStartNode.GetID()
			nextEnd := segEndNode.GetID()
			if prevStart.String() != nextStart.String() && prevStart.String() != nextEnd.String() &&
				prevEnd.String() != nextStart.String() && prevEnd.String() != nextEnd.String() {
				return nil, fmt.Errorf("failed to generate PathWrapper, Path received is invalid")
			}
		}
		segList = append(segList, segment{
			startNode:    segStartNode,
			relationship: relationship,
			endNode:      segEndNode,
		})
		src = dst
	}
	return &PathWrapper{
		path:             path,
		nodeList:         nodeList,
		relationshipList: relationshipList,
		segments:         segList,
		timezoneInfo:     timezoneInfo,
	}, nil
}

// AsStringTable Returns a 2D array of strings representing the query result
// If resultSet.resp.data is nil, returns an empty 2D array
func (res ResultSet) AsStringTable() [][]string {
	var resTable [][]string
	colNames := res.GetColNames()
	resTable = append(resTable, colNames)
	rows := res.GetRows()
	for _, row := range rows {
		var tempRow []string
		for _, val := range row.Values {
			tempRow = append(tempRow, ValueWrapper{val, res.timezoneInfo}.String())
		}
		resTable = append(resTable, tempRow)
	}
	return resTable
}

// GetValuesByColName Returns all values in the given column
func (res ResultSet) GetValuesByColName(colName string) ([]*ValueWrapper, error) {
	if !res.hasColName(colName) {
		return nil, fmt.Errorf("failed to get values, given column name '%s' does not exist", colName)
	}
	// Get index
	index := res.colNameIndexMap[colName]
	var valList []*ValueWrapper
	for _, row := range res.resp.Data.Rows {
		valList = append(valList, &ValueWrapper{row.Values[index], res.timezoneInfo})
	}
	return valList, nil
}

// GetRowValuesByIndex Returns all values in the row at given index
func (res ResultSet) GetRowValuesByIndex(index int) (*Record, error) {
	if err := checkIndex(index, res.resp.Data.Rows); err != nil {
		return nil, err
	}
	valWrap, err := genValWraps(res.resp.Data.Rows[index], res.timezoneInfo)
	if err != nil {
		return nil, err
	}
	return &Record{
		columnNames:     &res.columnNames,
		_record:         valWrap,
		colNameIndexMap: &res.colNameIndexMap,
		timezoneInfo:    res.timezoneInfo,
	}, nil
}

// Scan scans the rows into the given value.
func (res ResultSet) Scan(v interface{}) error {
	size := res.GetRowSize()
	if size == 0 {
		return nil
	}

	rv := reflect.ValueOf(v)
	switch {
	case rv.Kind() != reflect.Ptr:
		if t := reflect.TypeOf(v); t != nil {
			return fmt.Errorf("scan: Scan(non-pointer %s)", t)
		}
		fallthrough
	case rv.IsNil():
		return fmt.Errorf("scan: Scan(nil)")
	}
	rv = reflect.Indirect(rv)
	if k := rv.Kind(); k != reflect.Slice {
		return fmt.Errorf("scan: invalid type %s. expected slice as an argument", k)
	}

	colNames := res.GetColNames()
	rows := res.GetRows()

	t := reflect.TypeOf(v).Elem().Elem()
	for _, row := range rows {
		vv, err := res.scanRow(row, colNames, t)
		if err != nil {
			return err
		}
		rv.Set(reflect.Append(rv, vv))
	}

	return nil
}

// Scan scans the rows into the given value.
func (res ResultSet) scanRow(row *nebula.Row, colNames []string, rowType reflect.Type) (reflect.Value, error) {
	rowVals := row.GetValues()

	var result reflect.Value
	if rowType.Kind() == reflect.Ptr {
		result = reflect.New(rowType.Elem())
	} else {
		result = reflect.New(rowType).Elem()
	}
	structVal := reflect.Indirect(result)

	for fIdx := 0; fIdx < structVal.Type().NumField(); fIdx++ {
		f := structVal.Type().Field(fIdx)
		tag := f.Tag.Get("nebula")

		if tag == "" {
			continue
		}

		cIdx := slices.Index(colNames, tag)
		if cIdx == -1 {
			// It is possible that the tag is not in the result set
			continue
		}

		rowVal := rowVals[cIdx]

		if f.Type.Kind() == reflect.Slice {
			list := rowVal.GetLVal()
			err := scanListCol(list.Values, structVal.Field(fIdx), f.Type)
			if err != nil {
				return result, err
			}
		} else {
			err := scanPrimitiveCol(rowVal, structVal.Field(fIdx), f.Type.Kind())
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

func scanListCol(vals []*nebula.Value, listVal reflect.Value, sliceType reflect.Type) error {
	switch sliceType.Elem().Kind() {
	case reflect.Struct:
		var listCol = reflect.MakeSlice(sliceType, 0, len(vals))
		for _, val := range vals {
			ele := reflect.New(sliceType.Elem()).Elem()
			err := scanStructField(val, ele, sliceType.Elem())
			if err != nil {
				return err
			}
			listCol = reflect.Append(listCol, ele)
		}
		listVal.Set(listCol)
	case reflect.Ptr:
		var listCol = reflect.MakeSlice(sliceType, 0, len(vals))
		for _, val := range vals {
			ele := reflect.New(sliceType.Elem().Elem())
			err := scanStructField(val, reflect.Indirect(ele), sliceType.Elem().Elem())
			if err != nil {
				return err
			}
			listCol = reflect.Append(listCol, ele)
		}
		listVal.Set(listCol)
	default:
		return errors.New("scan: not support list type")
	}

	return nil
}

func scanStructField(val *nebula.Value, eleVal reflect.Value, eleType reflect.Type) error {
	vertex := val.GetVVal()
	if vertex != nil {
		tags := vertex.GetTags()
		vid := vertex.GetVid()

		if len(tags) != 0 {
			tag := tags[0]

			props := tag.GetProps()
			props["_vid"] = vid
			tagName := tag.GetName()
			props["_tag_name"] = &nebula.Value{SVal: tagName}

			err := scanValFromProps(props, eleVal, eleType)
			if err != nil {
				return err
			}
			return nil
		}
		// no tags, continue
	}

	edge := val.GetEVal()
	if edge != nil {
		props := edge.GetProps()

		src := edge.GetSrc()
		dst := edge.GetDst()
		name := edge.GetName()
		props["_src"] = src
		props["_dst"] = dst
		props["_name"] = &nebula.Value{SVal: name}

		err := scanValFromProps(props, eleVal, eleType)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func scanValFromProps(props map[string]*nebula.Value, val reflect.Value, tpe reflect.Type) error {
	for fIdx := 0; fIdx < tpe.NumField(); fIdx++ {
		f := tpe.Field(fIdx)
		n := f.Tag.Get("nebula")
		v, ok := props[n]
		if !ok {
			continue
		}
		err := scanPrimitiveCol(v, val.Field(fIdx), f.Type.Kind())
		if err != nil {
			return err
		}
	}

	return nil
}

func scanPrimitiveCol(rowVal *nebula.Value, val reflect.Value, kind reflect.Kind) error {
	w := ValueWrapper{value: rowVal}
	if w.IsNull() || w.IsEmpty() {
		// SetZero is introduced in go 1.20
		// val.SetZero()
		return nil
	}

	switch kind {
	case reflect.Bool:
		val.SetBool(rowVal.GetBVal())
	case reflect.Int:
		val.SetInt(rowVal.GetIVal())
	case reflect.Int8:
		val.SetInt(rowVal.GetIVal())
	case reflect.Int16:
		val.SetInt(rowVal.GetIVal())
	case reflect.Int32:
		val.SetInt(rowVal.GetIVal())
	case reflect.Int64:
		val.SetInt(rowVal.GetIVal())
	case reflect.Float32:
		val.SetFloat(rowVal.GetFVal())
	case reflect.Float64:
		val.SetFloat(rowVal.GetFVal())
	case reflect.String:
		val.SetString(string(rowVal.GetSVal()))
	default:
		return errors.New("scan: not support primitive type")
	}

	return nil
}

// GetRowSize Returns the number of total rows
func (res ResultSet) GetRowSize() int {
	if res.resp.Data == nil {
		return 0
	}
	return len(res.resp.Data.Rows)
}

// GetColSize Returns the number of total columns
func (res ResultSet) GetColSize() int {
	if res.resp.Data == nil {
		return 0
	}
	return len(res.resp.Data.ColumnNames)
}

// GetRows Returns all rows
func (res ResultSet) GetRows() []*nebula.Row {
	if res.resp.Data == nil {
		var empty []*nebula.Row
		return empty
	}
	return res.resp.Data.Rows
}

func (res ResultSet) GetColNames() []string {
	return res.columnNames
}

// GetErrorCode Returns an integer representing an error type
// 0    ErrorCode_SUCCEEDED
// -1   ErrorCode_E_DISCONNECTED
// -2   ErrorCode_E_FAIL_TO_CONNECT
// -3   ErrorCode_E_RPC_FAILURE
// -4   ErrorCode_E_BAD_USERNAME_PASSWORD
// -5   ErrorCode_E_SESSION_INVALID
// -6   ErrorCode_E_SESSION_TIMEOUT
// -7   ErrorCode_E_SYNTAX_ERROR
// -8   ErrorCode_E_EXECUTION_ERROR
// -9   ErrorCode_E_STATEMENT_EMPTY
// -10  ErrorCode_E_USER_NOT_FOUND
// -11  ErrorCode_E_BAD_PERMISSION
// -12  ErrorCode_E_SEMANTIC_ERROR
func (res ResultSet) GetErrorCode() ErrorCode {
	return ErrorCode(res.resp.ErrorCode)
}

func (res ResultSet) GetLatency() int64 {
	return res.resp.LatencyInUs
}

func (res ResultSet) GetLatencyInMs() int64 {
	return res.resp.LatencyInUs / 1000
}

func (res ResultSet) GetSpaceName() string {
	if res.resp.SpaceName == nil {
		return ""
	}
	return string(res.resp.SpaceName)
}

func (res ResultSet) GetErrorMsg() string {
	if res.resp.ErrorMsg == nil {
		return ""
	}
	return string(res.resp.ErrorMsg)
}

func (res ResultSet) IsSetPlanDesc() bool {
	return res.resp.PlanDesc != nil
}

func (res ResultSet) GetPlanDesc() *graph.PlanDescription {
	return res.resp.PlanDesc
}

func (res ResultSet) IsSetComment() bool {
	return res.resp.Comment != nil
}

func (res ResultSet) GetComment() string {
	if res.resp.Comment == nil {
		return ""
	}
	return string(res.resp.Comment)
}

func (res ResultSet) IsSetData() bool {
	return res.resp.Data != nil
}

func (res ResultSet) IsEmpty() bool {
	if !res.IsSetData() || len(res.resp.Data.Rows) == 0 {
		return true
	}
	return false
}

func (res ResultSet) IsSucceed() bool {
	return res.GetErrorCode() == ErrorCode_SUCCEEDED
}

func (res ResultSet) IsPartialSucceed() bool {
	return res.GetErrorCode() == ErrorCode_E_PARTIAL_SUCCEEDED
}

func (res ResultSet) hasColName(colName string) bool {
	if _, ok := res.colNameIndexMap[colName]; ok {
		return true
	}
	return false
}

// GetValueByIndex Returns value in the record at given column index
func (record Record) GetValueByIndex(index int) (*ValueWrapper, error) {
	if err := checkIndex(index, record._record); err != nil {
		return nil, err
	}
	return record._record[index], nil
}

// GetValueByColName Returns value in the record at given column name
func (record Record) GetValueByColName(colName string) (*ValueWrapper, error) {
	if !record.hasColName(colName) {
		return nil, fmt.Errorf("failed to get values, given column name '%s' does not exist", colName)
	}
	// Get index
	index := (*record.colNameIndexMap)[colName]
	return record._record[index], nil
}

func (record Record) String() string {
	var strList []string
	for _, val := range record._record {
		strList = append(strList, val.String())
	}
	return strings.Join(strList, ", ")
}

func (record Record) hasColName(colName string) bool {
	if _, ok := (*record.colNameIndexMap)[colName]; ok {
		return true
	}
	return false
}

// getRawID returns a list of row vid
func (node Node) getRawID() *nebula.Value {
	return node.vertex.GetVid()
}

// GetID returns a list of vid of node
func (node Node) GetID() ValueWrapper {
	return ValueWrapper{node.vertex.GetVid(), node.timezoneInfo}
}

// GetTags returns a list of tag names of node
func (node Node) GetTags() []string {
	return node.tags
}

// HasTag checks if node contains given label
func (node Node) HasTag(label string) bool {
	if _, ok := node.tagNameIndexMap[label]; ok {
		return true
	}
	return false
}

// Properties returns all properties of a tag
func (node Node) Properties(tagName string) (map[string]*ValueWrapper, error) {
	kvMap := make(map[string]*ValueWrapper)
	// Check if label exists
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exist in the Node", tagName)
	}
	index := node.tagNameIndexMap[tagName]
	for k, v := range node.vertex.Tags[index].Props {
		kvMap[k] = &ValueWrapper{v, node.timezoneInfo}
	}
	return kvMap, nil
}

// Keys returns all prop names of the given tag name
func (node Node) Keys(tagName string) ([]string, error) {
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exist in the Node", tagName)
	}
	var propNameList []string
	index := node.tagNameIndexMap[tagName]
	for k := range node.vertex.Tags[index].Props {
		propNameList = append(propNameList, k)
	}
	return propNameList, nil
}

// Values returns all prop values of the given tag name
func (node Node) Values(tagName string) ([]*ValueWrapper, error) {
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exist in the Node", tagName)
	}
	var propValList []*ValueWrapper
	index := node.tagNameIndexMap[tagName]
	for _, v := range node.vertex.Tags[index].Props {
		propValList = append(propValList, &ValueWrapper{v, node.timezoneInfo})
	}
	return propValList, nil
}

// String returns a string representing node
// Node format: ("VertexID" :tag1{k0: v0,k1: v1}:tag2{k2: v2})
func (node Node) String() string {
	var keyList []string
	var kvStr []string
	var tagStr []string
	vertex := node.vertex
	vid := vertex.GetVid()
	for _, tag := range vertex.GetTags() {
		kvs := tag.GetProps()
		tagName := tag.GetName()
		for k := range kvs {
			keyList = append(keyList, k)
		}
		sort.Strings(keyList)
		for _, k := range keyList {
			kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{kvs[k], node.timezoneInfo}.String())
			kvStr = append(kvStr, kvTemp)
		}
		tagStr = append(tagStr, fmt.Sprintf("%s{%s}", tagName, strings.Join(kvStr, ", ")))
		keyList = nil
		kvStr = nil
	}
	if len(tagStr) == 0 { // No tag
		return fmt.Sprintf("(%s)", ValueWrapper{vid, node.timezoneInfo}.String())
	}
	return fmt.Sprintf("(%s :%s)", ValueWrapper{vid, node.timezoneInfo}.String(), strings.Join(tagStr, " :"))
}

// IsEqualTo Returns true if two nodes have same vid
func (n1 Node) IsEqualTo(n2 *Node) bool {
	if n1.GetID().IsString() && n2.GetID().IsString() {
		s1, _ := n1.GetID().AsString()
		s2, _ := n2.GetID().AsString()
		return s1 == s2
	} else if n1.GetID().IsInt() && n2.GetID().IsInt() {
		s1, _ := n1.GetID().AsInt()
		s2, _ := n2.GetID().AsInt()
		return s1 == s2
	}
	return false
}

func (relationship Relationship) GetSrcVertexID() ValueWrapper {
	if relationship.edge.Type > 0 {
		return ValueWrapper{relationship.edge.GetSrc(), relationship.timezoneInfo}
	}
	return ValueWrapper{relationship.edge.GetDst(), relationship.timezoneInfo}
}

func (relationship Relationship) GetDstVertexID() ValueWrapper {
	if relationship.edge.Type > 0 {
		return ValueWrapper{relationship.edge.GetDst(), relationship.timezoneInfo}
	}
	return ValueWrapper{relationship.edge.GetSrc(), relationship.timezoneInfo}
}

func (relationship Relationship) GetEdgeName() string {
	return string(relationship.edge.Name)
}

func (relationship Relationship) GetRanking() int64 {
	return int64(relationship.edge.Ranking)
}

// Properties returns a map where the key is property name and the value is property name
func (relationship Relationship) Properties() map[string]*ValueWrapper {
	kvMap := make(map[string]*ValueWrapper)
	var (
		keyList   []string
		valueList []*ValueWrapper
	)
	for k, v := range relationship.edge.Props {
		keyList = append(keyList, k)
		valueList = append(valueList, &ValueWrapper{v, relationship.timezoneInfo})
	}

	for i := 0; i < len(keyList); i++ {
		kvMap[keyList[i]] = valueList[i]
	}
	return kvMap
}

// Keys returns a list of keys
func (relationship Relationship) Keys() []string {
	var keys []string
	for key := range relationship.edge.GetProps() {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a list of values wrapped as ValueWrappers
func (relationship Relationship) Values() []*ValueWrapper {
	var values []*ValueWrapper
	for _, value := range relationship.edge.GetProps() {
		values = append(values, &ValueWrapper{value, relationship.timezoneInfo})
	}
	return values
}

// String returns a string representing relationship
// Relationship format: [:edge src->dst @ranking {props}]
func (relationship Relationship) String() string {
	edge := relationship.edge
	var keyList []string
	var kvStr []string
	var src string
	var dst string
	for k := range edge.Props {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{edge.Props[k], relationship.timezoneInfo}.String())
		kvStr = append(kvStr, kvTemp)
	}
	if relationship.edge.Type > 0 {
		src = ValueWrapper{edge.Src, relationship.timezoneInfo}.String()
		dst = ValueWrapper{edge.Dst, relationship.timezoneInfo}.String()
	} else {
		src = ValueWrapper{edge.Dst, relationship.timezoneInfo}.String()
		dst = ValueWrapper{edge.Src, relationship.timezoneInfo}.String()
	}
	return fmt.Sprintf(`[:%s %s->%s @%d {%s}]`,
		string(edge.Name), src, dst, edge.Ranking, strings.Join(kvStr, ", "))
}

func (r1 Relationship) IsEqualTo(r2 *Relationship) bool {
	if r1.edge.GetSrc().IsSetSVal() && r2.edge.GetSrc().IsSetSVal() &&
		r1.edge.GetDst().IsSetSVal() && r2.edge.GetDst().IsSetSVal() {
		s1, _ := ValueWrapper{r1.edge.GetSrc(), r1.timezoneInfo}.AsString()
		s2, _ := ValueWrapper{r2.edge.GetSrc(), r2.timezoneInfo}.AsString()
		return s1 == s2 && string(r1.edge.Name) == string(r2.edge.Name) && r1.edge.Ranking == r2.edge.Ranking
	} else if r1.edge.GetSrc().IsSetIVal() && r2.edge.GetSrc().IsSetIVal() &&
		r1.edge.GetDst().IsSetIVal() && r2.edge.GetDst().IsSetIVal() {
		s1, _ := ValueWrapper{r1.edge.GetSrc(), r1.timezoneInfo}.AsInt()
		s2, _ := ValueWrapper{r2.edge.GetSrc(), r2.timezoneInfo}.AsInt()
		return s1 == s2 && string(r1.edge.Name) == string(r2.edge.Name) && r1.edge.Ranking == r2.edge.Ranking
	}
	return false
}

func (path *PathWrapper) GetPathLength() int {
	return len(path.segments)
}

func (path *PathWrapper) GetNodes() []*Node {
	return path.nodeList
}

func (path *PathWrapper) GetRelationships() []*Relationship {
	return path.relationshipList
}

func (path *PathWrapper) GetSegments() []segment {
	return path.segments
}

func (path *PathWrapper) ContainsNode(node Node) bool {
	for _, n := range path.nodeList {
		if n.IsEqualTo(&node) {
			return true
		}
	}
	return false
}

func (path *PathWrapper) ContainsRelationship(relationship *Relationship) bool {
	for _, r := range path.relationshipList {
		if r.IsEqualTo(relationship) {
			return true
		}
	}
	return false
}

func (path *PathWrapper) GetStartNode() (*Node, error) {
	if len(path.segments) == 0 {
		return nil, fmt.Errorf("failed to get start node, no node in the path")
	}
	return path.segments[0].startNode, nil
}

func (path *PathWrapper) GetEndNode() (*Node, error) {
	if len(path.segments) == 0 {
		return nil, fmt.Errorf("failed to get end node, no node in the path")
	}
	return path.segments[len(path.segments)-1].endNode, nil
}

// Path format: <("VertexID" :tag1{k0: v0,k1: v1})
// -[:TypeName@ranking {edgeProps}]->
// ("VertexID2" :tag1{k0: v0,k1: v1} :tag2{k2: v2})
// -[:TypeName@ranking {edgeProps}]->
// ("VertexID3" :tag1{k0: v0,k1: v1})>
func (pathWrap *PathWrapper) String() string {
	path := pathWrap.path
	src := path.Src
	steps := path.Steps
	resStr := ValueWrapper{&nebula.Value{VVal: src}, pathWrap.timezoneInfo}.String()
	for _, step := range steps {
		var keyList []string
		var kvStr []string
		for k := range step.Props {
			keyList = append(keyList, k)
		}
		sort.Strings(keyList)
		for _, k := range keyList {
			kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{step.Props[k], pathWrap.timezoneInfo}.String())
			kvStr = append(kvStr, kvTemp)
		}
		var dirChar1 string
		var dirChar2 string
		if step.Type > 0 {
			dirChar1 = "-"
			dirChar2 = "->"
		} else {
			dirChar1 = "<-"
			dirChar2 = "-"
		}
		resStr = resStr + fmt.Sprintf("%s[:%s@%d {%s}]%s%s",
			dirChar1,
			string(step.Name),
			step.Ranking,
			strings.Join(kvStr, ", "),
			dirChar2,
			ValueWrapper{&nebula.Value{VVal: step.Dst}, pathWrap.timezoneInfo}.String())
	}
	return "<" + resStr + ">"
}

func (p1 *PathWrapper) IsEqualTo(p2 *PathWrapper) bool {
	// Check length
	if len(p1.nodeList) != len(p2.nodeList) || len(p1.relationshipList) != len(p2.relationshipList) ||
		len(p1.segments) != len(p2.segments) {
		return false
	}
	// Check nodes
	for i := 0; i < len(p1.nodeList); i++ {
		if !p1.nodeList[i].IsEqualTo(p2.nodeList[i]) {
			return false
		}
	}
	// Check relationships
	for i := 0; i < len(p1.relationshipList); i++ {
		if !p1.relationshipList[i].IsEqualTo(p2.relationshipList[i]) {
			return false
		}
	}
	// Check segments
	for i := 0; i < len(p1.segments); i++ {
		if !p1.segments[i].startNode.IsEqualTo(p2.segments[i].startNode) ||
			!p1.segments[i].endNode.IsEqualTo(p2.segments[i].endNode) ||
			!p1.segments[i].relationship.IsEqualTo(p2.segments[i].relationship) {
			return false
		}
	}
	return true
}

func genTimeWrapper(time *nebula.Time, timezoneInfo timezoneInfo) (*TimeWrapper, error) {
	if time == nil {
		return nil, fmt.Errorf("failed to generate Time: invalid Time")
	}

	return &TimeWrapper{
		time:         time,
		timezoneInfo: timezoneInfo,
	}, nil
}

// getHour returns the hour in UTC
func (t TimeWrapper) getHour() int8 {
	return t.time.Hour
}

// getHour returns the minute in UTC
func (t TimeWrapper) getMinute() int8 {
	return t.time.Minute
}

// getHour returns the second in UTC
func (t TimeWrapper) getSecond() int8 {
	return t.time.Sec
}

func (t TimeWrapper) getMicrosec() int32 {
	return t.time.Microsec
}

// getRawTime returns a nebula.Time object in UTC.
//
//nolint:unused
func (t TimeWrapper) getRawTime() *nebula.Time {
	return t.time
}

// getLocalTime returns a nebula.Time object representing
// local time using timezone offset from the server.
func (t TimeWrapper) getLocalTime() (*nebula.Time, error) {
	// Original time object generated from server in UTC
	// Year, month and day are mocked up to fill the parameters
	rawTime := time.Date(2020,
		time.Month(1),
		1,
		int(t.getHour()),
		int(t.getMinute()),
		int(t.getSecond()),
		int(t.getMicrosec()*1000),
		time.UTC)

	// Use offset in seconds
	offset, err := time.ParseDuration(fmt.Sprintf("%ds", t.timezoneInfo.offset))
	if err != nil {
		return nil, err
	}
	localTime := rawTime.Add(offset)
	return &nebula.Time{
		Hour:     int8(localTime.Hour()),
		Minute:   int8(localTime.Minute()),
		Sec:      int8(localTime.Second()),
		Microsec: int32(localTime.Nanosecond() / 1000)}, nil
}

// getLocalTimeWithTimezoneOffset returns a nebula.Time object representing
// local time using user specified offset.
// Year, month, day in time.Time are filled with dummy values.
// Offset is in seconds.
func (t TimeWrapper) getLocalTimeWithTimezoneOffset(timezoneOffsetSeconds int32) (*nebula.Time, error) {
	// Original time object generated from server in UTC
	// Year, month and day are mocked up to fill the parameters
	rawTime := time.Date(2020,
		time.Month(1),
		1,
		int(t.getHour()),
		int(t.getMinute()),
		int(t.getSecond()),
		int(t.getMicrosec()*1000),
		time.UTC)

	offset, err := time.ParseDuration(fmt.Sprintf("%ds", timezoneOffsetSeconds))
	if err != nil {
		return nil, err
	}
	localTime := rawTime.Add(offset)
	return &nebula.Time{
		Hour:     int8(localTime.Hour()),
		Minute:   int8(localTime.Minute()),
		Sec:      int8(localTime.Second()),
		Microsec: int32(localTime.Nanosecond() / 1000)}, nil
}

// getLocalTimeWithTimezoneName returns a nebula.Time object
// representing local time using user specified timezone name.
// Year, month, day in time.Time are filled with 0.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
func (t TimeWrapper) getLocalTimeWithTimezoneName(timezoneName string) (*nebula.Time, error) {
	// Original time object generated from server in UTC
	// Year, month and day are mocked up to fill the parameters
	rawTime := time.Date(2020,
		time.Month(1),
		1,
		int(t.getHour()),
		int(t.getMinute()),
		int(t.getSecond()),
		int(t.getMicrosec()*1000),
		time.UTC)

	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return nil, err
	}
	localTime := rawTime.In(location)
	return &nebula.Time{
		Hour:     int8(localTime.Hour()),
		Minute:   int8(localTime.Minute()),
		Sec:      int8(localTime.Second()),
		Microsec: int32(localTime.Nanosecond() / 1000)}, nil
}

func (t1 TimeWrapper) IsEqualTo(t2 TimeWrapper) bool {
	return t1.getHour() == t2.getHour() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getMicrosec() == t2.getMicrosec()
}

func genDateWrapper(date *nebula.Date) (*DateWrapper, error) {
	if date == nil {
		return nil, fmt.Errorf("failed to generate date: invalid date")
	}
	return &DateWrapper{
		date: date,
	}, nil
}

func (d DateWrapper) getYear() int16 {
	return d.date.Year
}

func (d DateWrapper) getMonth() int8 {
	return d.date.Month
}

func (d DateWrapper) getDay() int8 {
	return d.date.Day
}

// getRawDate returns a nebula.Date object in UTC.
//
//nolint:unused
func (d DateWrapper) getRawDate() *nebula.Date {
	return d.date
}

func (d1 DateWrapper) IsEqualTo(d2 DateWrapper) bool {
	return d1.getYear() == d2.getYear() &&
		d1.getMonth() == d2.getMonth() &&
		d1.getDay() == d2.getDay()
}

func genDateTimeWrapper(datetime *nebula.DateTime, timezoneInfo timezoneInfo) (*DateTimeWrapper, error) {
	if datetime == nil {
		return nil, fmt.Errorf("failed to generate datetime: invalid datetime")
	}
	return &DateTimeWrapper{
		dateTime:     datetime,
		timezoneInfo: timezoneInfo,
	}, nil
}

func (dt DateTimeWrapper) getYear() int16 {
	return dt.dateTime.Year
}

func (dt DateTimeWrapper) getMonth() int8 {
	return dt.dateTime.Month
}

func (dt DateTimeWrapper) getDay() int8 {
	return dt.dateTime.Day
}

func (dt DateTimeWrapper) getHour() int8 {
	return dt.dateTime.Hour
}

func (dt DateTimeWrapper) getMinute() int8 {
	return dt.dateTime.Minute
}

func (dt DateTimeWrapper) getSecond() int8 {
	return dt.dateTime.Sec
}

func (dt DateTimeWrapper) getMicrosec() int32 {
	return dt.dateTime.Microsec
}

func (dt1 DateTimeWrapper) IsEqualTo(dt2 DateTimeWrapper) bool {
	return dt1.getYear() == dt2.getYear() &&
		dt1.getMonth() == dt2.getMonth() &&
		dt1.getDay() == dt2.getDay() &&
		dt1.getHour() == dt2.getHour() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getMicrosec() == dt2.getMicrosec()
}

// getRawDateTime returns a nebula.DateTime object representing local dateTime in UTC.
//
//nolint:unused
func (dt DateTimeWrapper) getRawDateTime() *nebula.DateTime {
	return dt.dateTime
}

// getLocalDateTime returns a nebula.DateTime object representing
// local datetime using timezone offset from the server.
func (dt DateTimeWrapper) getLocalDateTime() (*nebula.DateTime, error) {
	// Original time object generated from server in UTC
	rawTime := time.Date(
		int(dt.getYear()), time.Month(dt.getMonth()), int(dt.getDay()),
		int(dt.getHour()), int(dt.getMinute()), int(dt.getSecond()), int(dt.dateTime.Microsec*1000),
		time.UTC)

	// Use offset in seconds
	offset, err := time.ParseDuration(fmt.Sprintf("%ds", dt.timezoneInfo.offset))
	if err != nil {
		return nil, err
	}
	localDT := rawTime.Add(offset)
	return &nebula.DateTime{
		Year:     int16(localDT.Year()),
		Month:    int8(localDT.Month()),
		Day:      int8(localDT.Day()),
		Hour:     int8(localDT.Hour()),
		Minute:   int8(localDT.Minute()),
		Sec:      int8(localDT.Second()),
		Microsec: int32(localDT.Nanosecond() / 1000)}, nil
}

// getLocalDateTimeWithTimezoneOffset returns a nebula.DateTime object representing
// local datetime using user specified timezone offset.
// Offset is in seconds.
func (dt DateTimeWrapper) getLocalDateTimeWithTimezoneOffset(timezoneOffsetSeconds int32) (*nebula.DateTime, error) {
	// Original time object generated from server in UTC
	rawTime := time.Date(
		int(dt.getYear()), time.Month(dt.getMonth()), int(dt.getDay()),
		int(dt.getHour()), int(dt.getMinute()), int(dt.getSecond()), int(dt.dateTime.Microsec*1000),
		time.UTC)

	offset, err := time.ParseDuration(fmt.Sprintf("%ds", timezoneOffsetSeconds))
	if err != nil {
		return nil, err
	}
	localDT := rawTime.Add(offset)
	return &nebula.DateTime{
		Year:     int16(localDT.Year()),
		Month:    int8(localDT.Month()),
		Day:      int8(localDT.Day()),
		Hour:     int8(localDT.Hour()),
		Minute:   int8(localDT.Minute()),
		Sec:      int8(localDT.Second()),
		Microsec: int32(localDT.Nanosecond() / 1000)}, nil
}

// GetLocalDateTimeWithTimezoneName returns a nebula.DateTime object representing
// local time using user specified timezone name.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
func (dt DateTimeWrapper) GetLocalDateTimeWithTimezoneName(timezoneName string) (*nebula.DateTime, error) {
	// Original time object generated from server in UTC
	rawTime := time.Date(
		int(dt.getYear()), time.Month(dt.getMonth()), int(dt.getDay()),
		int(dt.getHour()), int(dt.getMinute()), int(dt.getSecond()), int(dt.getMicrosec()*1000),
		time.UTC)

	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return nil, err
	}
	localDT := rawTime.In(location)

	return &nebula.DateTime{
		Year:     int16(localDT.Year()),
		Month:    int8(localDT.Month()),
		Day:      int8(localDT.Day()),
		Hour:     int8(localDT.Hour()),
		Minute:   int8(localDT.Minute()),
		Sec:      int8(localDT.Second()),
		Microsec: int32(localDT.Nanosecond() / 1000)}, nil
}

func checkIndex(index int, list interface{}) error {
	if _, ok := list.([]*nebula.Row); ok {
		if index < 0 || index >= len(list.([]*nebula.Row)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	} else if _, ok := list.([]*ValueWrapper); ok {
		if index < 0 || index >= len(list.([]*ValueWrapper)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	}
	return fmt.Errorf("given list type is invalid")
}

func graphvizString(s string) string {
	s = strings.Replace(s, "{", "\\{", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "<", "\\<", -1)
	s = strings.Replace(s, ">", "\\>", -1)
	return s
}

func prettyFormatJsonString(value []byte) string {
	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, value, "", "  "); err != nil {
		return string(value)
	}
	return prettyJson.String()
}

func name(planNodeDesc *graph.PlanNodeDescription) string {
	return fmt.Sprintf("%s_%d", planNodeDesc.GetName(), planNodeDesc.GetID())
}

func condEdgeLabel(condNode *graph.PlanNodeDescription, doBranch bool) string {
	name := strings.ToLower(string(condNode.GetName()))
	if strings.HasPrefix(name, "select") {
		if doBranch {
			return "Y"
		}
		return "N"
	}
	if strings.HasPrefix(name, "loop") {
		if doBranch {
			return "Do"
		}
	}
	return ""
}

func nodeString(planNodeDesc *graph.PlanNodeDescription, planNodeName string) string {
	var outputVar = graphvizString(string(planNodeDesc.GetOutputVar()))
	var inputVar string
	if planNodeDesc.IsSetDescription() {
		desc := planNodeDesc.GetDescription()
		for _, pair := range desc {
			key := string(pair.GetKey())
			if key == "inputVar" {
				inputVar = graphvizString(string(pair.GetValue()))
			}
		}
	}
	return fmt.Sprintf("\t\"%s\"[label=\"{%s|outputVar: %s|inputVar: %s}\", shape=Mrecord];\n",
		planNodeName, planNodeName, outputVar, inputVar)
}

func edgeString(start, end string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\";\n", start, end)
}

func conditionalEdgeString(start, end, label string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\"[label=\"%s\", style=dashed];\n", start, end, label)
}

func conditionalNodeString(name string) string {
	return fmt.Sprintf("\t\"%s\"[shape=diamond];\n", name)
}

func nodeById(p *graph.PlanDescription, nodeId int64) *graph.PlanNodeDescription {
	line := p.GetNodeIndexMap()[nodeId]
	return p.GetPlanNodeDescs()[line]
}

func findBranchEndNode(p *graph.PlanDescription, condNodeId int64, isDoBranch bool) int64 {
	for _, node := range p.GetPlanNodeDescs() {
		if node.IsSetBranchInfo() {
			bInfo := node.GetBranchInfo()
			if bInfo.GetConditionNodeID() == condNodeId && bInfo.GetIsDoBranch() == isDoBranch {
				return node.GetID()
			}
		}
	}
	return -1
}

func findFirstStartNodeFrom(p *graph.PlanDescription, nodeId int64) int64 {
	node := nodeById(p, nodeId)
	for {
		deps := node.GetDependencies()
		if len(deps) == 0 {
			if strings.ToLower(string(node.GetName())) != "start" {
				return -1
			}
			return node.GetID()
		}
		node = nodeById(p, deps[0])
	}
}

// explain/profile format="dot"
func (res ResultSet) MakeDotGraph() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		case "select":
			builder.WriteString(conditionalNodeString(planNodeName))
			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
			// then branch
			thenNodeId := findBranchEndNode(p, planNodeDesc.GetID(), true)
			builder.WriteString(edgeString(name(nodeById(p, thenNodeId)), name(dep)))
			thenStartId := findFirstStartNodeFrom(p, thenNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, thenStartId)), "Y"))
			// else branch
			elseNodeId := findBranchEndNode(p, planNodeDesc.GetID(), false)
			builder.WriteString(edgeString(name(nodeById(p, elseNodeId)), name(dep)))
			elseStartId := findFirstStartNodeFrom(p, elseNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, elseStartId)), "N"))
			// dep
			builder.WriteString(edgeString(name(dep), planNodeName))
		case "loop":
			builder.WriteString(conditionalNodeString(planNodeName))
			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
			// do branch
			doNodeId := findBranchEndNode(p, planNodeDesc.GetID(), true)
			builder.WriteString(edgeString(name(nodeById(p, doNodeId)), name(planNodeDesc)))
			doStartId := findFirstStartNodeFrom(p, doNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, doStartId)), "Do"))
			// dep
			builder.WriteString(edgeString(name(dep), planNodeName))
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
			if planNodeDesc.IsSetDependencies() {
				for _, depId := range planNodeDesc.GetDependencies() {
					builder.WriteString(edgeString(name(nodeById(p, depId)), planNodeName))
				}
			}
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// explain/profile format="dot:struct"
func (res ResultSet) MakeDotGraphByStruct() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		case "select":
			builder.WriteString(conditionalNodeString(planNodeName))
		case "loop":
			builder.WriteString(conditionalNodeString(planNodeName))
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
		}

		if planNodeDesc.IsSetDependencies() {
			for _, depId := range planNodeDesc.GetDependencies() {
				dep := nodeById(p, depId)
				builder.WriteString(edgeString(name(dep), planNodeName))
			}
		}

		if planNodeDesc.IsSetBranchInfo() {
			branchInfo := planNodeDesc.GetBranchInfo()
			condNode := nodeById(p, branchInfo.GetConditionNodeID())
			label := condEdgeLabel(condNode, branchInfo.GetIsDoBranch())
			builder.WriteString(conditionalEdgeString(planNodeName, name(condNode), label))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// MakeProfilingData generate profiling data for both Row and TCK formats.
func MakeProfilingData(planNodeDesc *graph.PlanNodeDescription, isTckFmt bool) (string, error) {
	var profileArr []string
	re, err := regexp.Compile(`^[^{(\[]\w+`)
	if err != nil {
		panic(err)
	}
	for i, profile := range planNodeDesc.GetProfiles() {
		var statArr []string
		statArr = append(statArr, fmt.Sprintf("\"version\":%d", i))
		statArr = append(statArr, fmt.Sprintf("\"rows\":%d", profile.GetRows()))
		if !isTckFmt {
			// tck format doesn't need these fields
			statArr = append(statArr, fmt.Sprintf("\"execTime\":\"%d(us)\"", profile.GetExecDurationInUs()))
			statArr = append(statArr, fmt.Sprintf("\"totalTime\":\"%d(us)\"", profile.GetTotalDurationInUs()))
		}
		for k, v := range profile.GetOtherStats() {
			s := string(v)
			if matched := re.Match(v); matched {
				if !strings.HasPrefix(s, "\"") {
					s = fmt.Sprintf("\"%s", s)
				}
				if !strings.HasSuffix(s, "\"") {
					s = fmt.Sprintf("%s\"", s)
				}
			}
			statArr = append(statArr, fmt.Sprintf("\"%s\": %s", k, s))
		}
		sort.Strings(statArr)
		statStr := fmt.Sprintf("{%s}", strings.Join(statArr, ",\n"))
		profileArr = append(profileArr, statStr)
	}
	allProfiles := strings.Join(profileArr, ",\n")
	if len(profileArr) > 1 {
		allProfiles = fmt.Sprintf("[%s]", allProfiles)
	}
	var buffer bytes.Buffer
	err = json.Indent(&buffer, []byte(allProfiles), "", "  ")
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// generate operator info for Row format.
func MakeOperatorInfo(planNodeDesc *graph.PlanNodeDescription) string {
	var columnInfo []string
	if planNodeDesc.IsSetBranchInfo() {
		branchInfo := planNodeDesc.GetBranchInfo()
		columnInfo = append(columnInfo, fmt.Sprintf("branch: %t, nodeId: %d\n",
			branchInfo.GetIsDoBranch(), branchInfo.GetConditionNodeID()))
	}

	outputVar := fmt.Sprintf("outputVar: %s", prettyFormatJsonString(planNodeDesc.GetOutputVar()))
	columnInfo = append(columnInfo, outputVar)

	if planNodeDesc.IsSetDescription() {
		desc := planNodeDesc.GetDescription()
		for _, pair := range desc {
			value := prettyFormatJsonString(pair.GetValue())
			columnInfo = append(columnInfo, fmt.Sprintf("%s: %s", string(pair.GetKey()), value))
		}
	}
	return strings.Join(columnInfo, "\n")
}

// MakePlanByRow explain/profile format="row"
func (res ResultSet) MakePlanByRow() ([][]interface{}, error) {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var rows [][]interface{}
	for _, planNodeDesc := range planNodeDescs {
		var row []interface{}
		row = append(row, planNodeDesc.GetID(), string(planNodeDesc.GetName()))

		if planNodeDesc.IsSetDependencies() {
			var deps []string
			for _, dep := range planNodeDesc.GetDependencies() {
				deps = append(deps, fmt.Sprintf("%d", dep))
			}
			row = append(row, strings.Join(deps, ","))
		} else {
			row = append(row, "")
		}

		if planNodeDesc.IsSetProfiles() {
			r, err := MakeProfilingData(planNodeDesc, false)
			if err != nil {
				return nil, err
			}
			row = append(row, r)
		} else {
			row = append(row, "")
		}
		opInfo := MakeOperatorInfo(planNodeDesc)
		row = append(row, opInfo)
		rows = append(rows, row)
	}
	return rows, nil
}

// MakePlanByTck explain/profile format="tck"
func (res ResultSet) MakePlanByTck() ([][]interface{}, error) {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var rows [][]interface{}
	for _, planNodeDesc := range planNodeDescs {
		var row []interface{}
		row = append(row, planNodeDesc.GetID(), string(planNodeDesc.GetName()))

		if planNodeDesc.IsSetDependencies() {
			var deps []string
			for _, dep := range planNodeDesc.GetDependencies() {
				deps = append(deps, fmt.Sprintf("%d", dep))
			}
			row = append(row, strings.Join(deps, ","))
		} else {
			row = append(row, "")
		}

		if planNodeDesc.IsSetProfiles() {
			var compactProfilingData bytes.Buffer
			r, err := MakeProfilingData(planNodeDesc, true)
			if err != nil {
				return nil, err
			}
			// compress JSON data and remove whitespace characters
			err = json.Compact(&compactProfilingData, []byte(r))
			if err != nil {
				return nil, err
			}
			row = append(row, compactProfilingData.String())
		} else {
			row = append(row, "")
		}
		// append operator info
		row = append(row, "")

		rows = append(rows, row)
	}
	return rows, nil
}
