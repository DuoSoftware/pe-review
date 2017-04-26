// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package hive_metastore

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/mattbaird/hive/thriftlib/fb303"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

var _ = fb303.GoUnusedProtection__

const DDL_TIME = "transient_lastDdlTime"
const IS_ARCHIVED = "is_archived"
const ORIGINAL_LOCATION = "original_location"
const META_TABLE_COLUMNS = "columns"
const META_TABLE_COLUMN_TYPES = "columns.types"
const BUCKET_FIELD_NAME = "bucket_field_name"
const BUCKET_COUNT = "bucket_count"
const FIELD_TO_DIMENSION = "field_to_dimension"
const META_TABLE_NAME = "name"
const META_TABLE_DB = "db"
const META_TABLE_LOCATION = "location"
const META_TABLE_SERDE = "serde"
const META_TABLE_PARTITION_COLUMNS = "partition_columns"
const FILE_INPUT_FORMAT = "file.inputformat"
const FILE_OUTPUT_FORMAT = "file.outputformat"
const META_TABLE_STORAGE = "storage_handler"

func init() {
}
